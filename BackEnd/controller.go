package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type RichiestaLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GestioneLogin(db*gorm.DB) gin.HandlerFunc{
	return func (c*gin.Context){
		var req RichiestaLogin
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Dati non validi"})
			return 
		}
		var utente Utente 
		if err := db.Preload("Ruolo").Where("Username = ?", req.Username).First(&utente).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenziali errate"})
			return
		}
		if utente.StatoAccount == "Bloccato"{
			c.JSON(http.StatusForbidden, gin.H{"error":"Account bloccato per troppi tentativi falliti"})
		}
		if !CheckPasswordHash(req.Password, utente.Password){
			utente.TentativiFalliti++
			esito := "Fallito"

			if utente.TentativiFalliti >= 5 {
				utente.StatoAccount = "Bloccato"
				esito="Fallito, account bloccato"
			}
			db.Save(&utente)

			db.Create(&LogAccessi{
				UtenteID: utente.ID,
				DataAccesso: time.Now(),
				Esito: esito,
				IndirizzoIP: c.ClientIP(),
			})

			c.JSON(http.StatusUnauthorized, gin.H{"error":"Password errata"})
			return 
		}

		utente.TentativiFalliti = 0
		ora := time.Now()
		utente.UltimoLogin = &ora
		db.Save(&utente)

		db.Create(&LogAccessi{
			UtenteID: utente.ID,
			DataAccesso: time.Now(),
			Esito: "Successo",
			IndirizzoIP: c.ClientIP(),
		})

		access, refresh, err := GenerateTokens(utente.ID, utente.Ruolo.NomeRuolo)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Errore generazione token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": access,
			"refresh_token": refresh,
			"user":gin.H{
				"username":utente.Username,
				"ruolo":utente.Ruolo.NomeRuolo,
			},
		})
	}
}
func TokenCasuale() string{
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
func RichiestaResetPassword(db*gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"Email obbligatoria"})
			return
		}

		var utente Utente
		if err := db.Where("email = ?", req.Email).First(&utente).Error; err != nil{
			c.JSON(http.StatusOK, gin.H{"message":"Se l'email è valida riceverai un link"})
			return
		}
		token := TokenCasuale()
		scadenza := time.Now().Add(1*time.Hour)

		reset := RecuperoPassword{
			UtenteID: utente.ID,
			TokenUnivoco: token,
			DataScadenza: scadenza,
			Stato: "Non usato",
		}

		if err := db.Create(&reset).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Errore interno"})
			return 
		}

		// Simulazione invio email
		fmt.Printf("\n--- EMAIL DI RECUPERO per %s ---\nLink: http://localhost:3000/reset-password?token=%s\n--------------------------------\n", utente.Email, token)

		c.JSON(http.StatusOK, gin.H{"message": "Istruzioni inviate via email"})

	}
}
func ConfermaResetPassword(db * gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var req struct{
			Token       string `json:"token" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"Dati mancanti"})
			return 
		}
		if !PasswordValida(req.NewPassword){
			c.JSON(http.StatusBadRequest, gin.H{"error":"La password non rispetta i criteri AGID (min 8 caratteri, maiuscola, numero, speciale)"})
			return 
		}

		var reset RecuperoPassword

		err := db.Where("token_univoco = ? AND stato = ? AND data_scadenza > ?", req.Token, "Non usato", time.Now()).First(&reset).Error
		if err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Token non valido o scaduto"})
			return
		}

		hashedPassword, _ := HashPassword(req.NewPassword)


		err = db.Transaction(func(tx *gorm.DB) error {
			if err :=tx.Model(&Utente{}).Where("id = ?", reset.UtenteID).Update("password", hashedPassword).Error; err != nil{
				return err
			}

			if err := tx.Model(&reset).Update("stato", "Usato").Error; err != nil{
				return err
			}
			return nil
		})

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Errore durante il reset"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message":"Password aggiornata con successo"})

	}
}
type GraficoPerCategoria struct{
	Nome string `json:"nome"`
	Quantita string `json:"quantita"`
}
type StatisticheDashboard struct{
	UtentiTotali int64 `json:"utenti_totali"`
	ProdottiTotali int64 `json:"prodotti_totali"`
	ValoreInventario float64 `json:"valore_inventario"`
	UltimiMovimenti []MovimentoMagazzino `json:"ultimi_movimenti"`
	GraficoCategoria []GraficoPerCategoria `json:"grafico_categorie"`
	
}
func GetStatisticheDashboard(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {


		var statistiche StatisticheDashboard 
	
		db.Model(&Utente{}).Count(&statistiche.UtentiTotali)  //conto utenti totali
		db.Model(&Prodotto{}).Count(&statistiche.ProdottiTotali)  //conto prodotti totali
		
		db.Model(&Prodotto{}).Select("COALESCE(SUM(quantita_disponibile*prezzo_unitario),0)").Row().Scan(&statistiche.ValoreInventario) //calcolo prezzo dell'intero inventario, se è vuoto restituisce NULL
		
		db.Preload("Prodotto").Preload("UtenteOperazione").Order("data_movimento desc").Limit(5).Find(&statistiche.UltimiMovimenti) // Ultimi 5 movimenti

		db.Table("prodotto").Select("tipo_prodotto.corpo_messaggio as nome, count(prodotto.id) as quantita").Joins("join tipo_prodotto on tipo_prodotto.id = prodotto.tipo_prodotto_id").Group("tipo_prodotto.corpo_messaggio").Scan(&statistiche.GraficoCategoria)
	
		c.JSON(200, statistiche)
	}
	
}

func GestioneUtenti(db*gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		page, _:= strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _:= strconv.Atoi(c.DefaultQuery("limit", "10"))
		search := c.Query("search")
		offset  := (page - 1)* limit

		var utenti []Utente
		var total int64

		query := db.Model(&Utente{}).Preload("Ruolo")

		if search != "" {
			query = query.Where("username ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
		}

		query.Count(&total)

		if err := query.Limit(limit).Offset(offset).Find(&utenti).Error; err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Errore nel recupero utenti"})
			return 
		}

		c.JSON(http.StatusOK, gin.H{
			"data": utenti,
			"total": total,
			"page": page,
			"last_page": int(total/int64(limit))+1,
		})


	}
}

func DisattivaUtente(db*gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Where("id = ?", id).Delete(&Utente{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Errore nella disattivazione dell'utente"})
		}

		c.JSON(http.StatusOK, gin.H{"message":"Utente disattivato"})
	}
}


func CreaUtente(db*gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var req struct{
		Username    string    `json:"username" binding:"required"`
		Email       string    `json:"email" binding:"required"`
		Password    string    `json:"password" binding:"required"`
		Nome        string    `json:"nome"`
		Cognome     string    `json:"cognome"`
		DataNascita time.Time `json:"data_nascita"`
		RuoloID     uuid.UUID `json:"ruolo_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Dati mancanti o errati"})
		return 
	}

	if !PasswordValida(req.Password){
		c.JSON(http.StatusBadRequest, gin.H{"error":"La Password non rispetta le linee guida"})
		return
	}

	hashedPassword, _ := HashPassword(req.Password)

	NewUtente := Utente{
		Username:     req.Username,
		Email:        req.Email,
		Password:     hashedPassword,
		Nome:         req.Nome,
		Cognome:      req.Cognome,
		DataNascita:  req.DataNascita,
		RuoloID:      req.RuoloID,
		StatoAccount: "Attivo",
	}

	if err := db.Create(&NewUtente).Error; err != nil{
		c.JSON(http.StatusConflict, gin.H{"error":"Username o Email già esistenti"})
		return 
	}

	fmt.Printf("\n-- EMAIL DI BENVENUTO inviata a %s ---\nOggetto: Benvenuto nel Gestionale\nAccount creato con successo", NewUtente.Email)

	c.JSON(http.StatusCreated,NewUtente)

	}
	

}
func ModificaUtente(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id") 
		
		var req struct {
			Nome        string    `json:"nome"`
			Cognome     string    `json:"cognome"`
			Email       string    `json:"email"`
			RuoloID     uuid.UUID `json:"ruolo_id"`
			DataNascita time.Time `json:"data_nascita"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
			return
		}

		result := db.Model(&Utente{}).Where("id = ?", id).Updates(map[string]interface{}{
			"nome":         req.Nome,
			"cognome":      req.Cognome,
			"email":        req.Email,
			"ruolo_id":     req.RuoloID, 
			"data_nascita": req.DataNascita,
		})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante l'aggiornamento"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Utente e ruolo aggiornati correttamente"})
	}
}
