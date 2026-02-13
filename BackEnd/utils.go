package main

import (
	"golang.org/x/crypto/bcrypt"
)

//trasforma la password da in chiaro in hash sicuro
func HashPassword(password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes), err
}

//confronta la password inserita con l'hash nel DB
func CheckPasswordHash(password, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}