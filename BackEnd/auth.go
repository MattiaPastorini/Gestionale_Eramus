package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte("JWT_SECRET")

type Claims struct {
	UtenteID uuid.UUID `json:"utente_id"`
	Ruolo string `json:"ruolo"`
	jwt.RegisteredClaims
}



func GenerateTokens (UtenteID uuid.UUID, Ruolo string) (string, string, error){
	expirationTime := time.Now().Add(15*time.Minute)
	claims := &Claims{
		UtenteID: UtenteID,
		Ruolo: Ruolo,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessString, err := accessToken.SignedString(jwtKey)


	refreshExpiration := time.Now().Add(7*24*time.Hour)
	refreshClaims := &Claims{
		UtenteID: UtenteID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(jwtKey)

	return accessString, refreshString, err
}