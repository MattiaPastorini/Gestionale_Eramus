package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


func main(){
	dsn := "host=localhost user=postgres password=1234 dbname=Eramus_Gestionale port=5432 sslmode=disable"
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