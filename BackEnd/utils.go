package main

import (
	"unicode"

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

func PasswordValida (s string) bool {
	var (
		hasMinLen = len(s) >= 8
		hasUpper = false
		hasNumber = false
		hasSpecial = false
	)
	for _, char := range s {
		switch{
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true	

		}
	}
	return hasMinLen && hasUpper && hasNumber && hasSpecial
}