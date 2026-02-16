package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ruoloRichiesto string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Accesso negato: token mancante"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil 
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessione scaduta o token non valido"})
			c.Abort()
			return
		}

		if ruoloRichiesto != "" && claims.Ruolo != ruoloRichiesto {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permessi insufficienti: solo " + ruoloRichiesto + " può farlo"})
			c.Abort()
			return
		}

		c.Set("utente_id", claims.UtenteID)
		c.Set("ruolo", claims.Ruolo)
		
		c.Next() 
	}
}

func LoggerMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if status >= 400{
			log.Printf("Allerta sicurezza \n%d \n%s \n%s \n%s \n%s",
			status,
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		
		)
		}
	}
}