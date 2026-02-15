package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)



func setupInitialData(db *gorm.DB) {
	var adminRuolo Ruolo
	// Crea ruolo se non esiste
	db.FirstOrCreate(&adminRuolo, Ruolo{NomeRuolo: "Admin"})
	db.FirstOrCreate(&Ruolo{}, Ruolo{NomeRuolo: "Operatore"})

	var count int64
	db.Model(&Utente{}).Count(&count)
	if count == 0 {
		// Usa la funzione HashPassword che hai in utils.go
		pass, _ := HashPassword("Admin123!") 
		admin := Utente{
			Username:     "admin",
			Password:     pass,
			StatoAccount: "Attivo",
			RuoloID:      adminRuolo.ID,
		}
		db.Create(&admin)
		fmt.Println("Utente Admin creato (Pass: Admin123!)")
	}
}



func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}

func main() {
    // Carica .env
    if err := godotenv.Load(); err != nil {
        log.Println("WARN: .env non trovato, uso variabili di sistema")
    }

    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=%s port=5432 sslmode=disable", 
        dbPassword, dbName)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
    })
    if err != nil {
        panic("Errore connessione al DataBase: " + err.Error())
    }

    // Config connection pool
    sqlDB, err := db.DB()
    if err != nil {
        panic("Errore DB pool: " + err.Error())
    }
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    // Migrazione
    if err := db.AutoMigrate(&Ruolo{}, &Utente{}, &LogAccessi{}); err != nil {
        log.Println("Errore migrazione:", err)
    } else {
        fmt.Println("Database sincronizzato")
    }

    setupInitialData(db)

    // Gin
    r := gin.Default()
    r.Use(CORSMiddleware())

    api := r.Group("/api")
    {
        api.POST("/login", GestioneLogin(db))
        api.POST("/forgot-password", RichiestaResetPassword(db))
        api.POST("/reset-password-confirm", ConfermaResetPassword(db))

		api.GET("/dashboard/statistiche", GetStatisticheDashboard(db))
    }

    fmt.Println("Server su http://localhost:8080")
    r.Run(":8080")
}
