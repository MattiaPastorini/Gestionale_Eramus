package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

		// Simulazione invio email (come richiesto, per ora in console)
		fmt.Printf("\n--- EMAIL DI RECUPERO per %s ---\nLink: http://localhost:3000/reset-password?token=%s\n--------------------------------\n", utente.Email, token)

		c.JSON(http.StatusOK, gin.H{"message": "Istruzioni inviate via email"})

	}
}
func ConfermaResetPassword (db * gorm.DB) gin.HandlerFunc{
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