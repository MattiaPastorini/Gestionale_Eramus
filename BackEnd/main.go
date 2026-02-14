package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


func main(){

	err:= godotenv.Load()
	if err != nil{
		log.Fatal("Errore nel caricamento del file .env")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=%s port=5432 sslmode=disable", dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // non aggiunge la "s" alle tabelle
		},
	})
	if err != nil{
		panic("Errore connessione al DataBase: " + err.Error())
	}

	err = db.AutoMigrate(
		&Ruolo{}, 
        &Utente{}, 
        &TipoProdotto{}, 
        &Prodotto{}, 
        &MovimentoMagazzino{}, 
        &LogAccessi{}, 
        &RecuperoPassword{}, 
        &NotificheEmail{},
	)

	if err != nil{
		fmt.Println("Errore migrazione: ", err)
	} else {
		fmt.Println("DataBase sincronizzato con successo")
	}
}