package main

import (
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
				UtenteID: utente.Id,
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
			UtenteID: utente.Id,
			DataAccesso: time.Now(),
			Esito: "Successo",
			IndirizzoIP: c.ClientIP(),
		})

		access, refresh, err := GenerateTokens(utente.Id, utente.Ruolo.NomeRuolo)
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