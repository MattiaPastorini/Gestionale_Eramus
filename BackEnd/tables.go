package main

import (
	"time"

	"github.com/google/uuid"
)

type Utente struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Nome string `gorm:"column:nome"`
	Cognome string `gorm:"column:cognome"`
	DataNascita time.Time `gorm:"column:data_nascita"`
	RuoloID uuid.UUID `gorm:"type:uuid;not null"`
	Ruolo Ruolo `gorm:"foreignKey:RuoloID"`	
	TentativiFalliti int `gorm:"default:0"`
	StatoAccount string `gorm:"default:'Attivo'"`
	UltimoLogin *time.Time `gorm:"column:data_ultimo_login"`
	DataCreazione time.Time `gorm:"column:data_creazione;autoCreateTime"`
	DataAggiornamento time.Time `gorm:"column:data_aggiornamento;autoUpdateTime"`
}
type Ruolo struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	NomeRuolo string `gorm:"unique;not null;column:nome_ruolo"`
	Descrizione string `gorm:"type:text"`
}
type TipoProdotto struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CorpoMessaggio string `gorm:"type:text"`
	DataInvio time.Time `gorm:"column:data_invio"`
	EsitoInvio string `gorm:"type:varchar(50)"`
}
type Prodotto struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	NomeOggetto string `gorm:"column:nome_oggetto"`
	Descrizione string `gorm:"type:text"`
	QuantitaDisponibile int `gorm:"default:0"`
	PrezzoUnitario float64 `gorm:"type:decimal(10,2);not null"`
	SogliaMinimaDiMagazzino int `gorm:"default:5"`
	DataInserimento time.Time `gorm:"column:data_inserimento;autoCreateTime"`
	TipoProdottoID uuid.UUID `gorm:"type:uuid;not null"`
	TipoProdotto TipoProdotto  `gorm:"foreignKey:TipoProdottoID"`	
	CreatoDaID uuid.UUID `gorm:"type:uuid;not null"`
	CreatoDa Utente `gorm:"foreignKey:CreatoDaID"`
	DataUltimaModifica time.Time `gorm:"column:data_ultima_modifica;autoUpdateTime"`
}
type MovimentoMagazzino struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProdottoID uuid.UUID `gorm:"type:uuid;not null"`
	Prodotto Prodotto  `gorm:"foreignKey:ProdottoID"`
	TipoMovimento string `gorm:"not null;column:tipo_movimento"`
	Quantita int `gorm:"default:0"`
	DataMovimento time.Time `gorm:"column:data_movimento;autoCreateTime"`
	UtenteOperazioneID uuid.UUID `gorm:"type:uuid;not null"`
	UtenteOperazione Utente `gorm:"foreignKey:UtenteOperazioneID"`
	Note string	`gorm:"type:text"`
}
type LogAccessi struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UtenteID uuid.UUID `gorm:"type:uuid;not null"`
	Utente Utente `gorm:"foreignKey:UtenteID"`
	DataAccesso time.Time `gorm:"column:data_accesso;autoCreateTime"`
	Esito string `gorm:"type:varchar(50)"`
	IndirizzoIP string `gorm:"type:varchar(45)"` 
}

type RecuperoPassword struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UtenteID uuid.UUID `gorm:"type:uuid;not null"`
	Utente Utente `gorm:"foreignKey:UtenteID"`
	TokenUnivoco string    `gorm:"unique;not null"`
	DataGenerazione  time.Time `gorm:"column:data_generazione;autoCreateTime"`
	DataScadenza time.Time `gorm:"column:data_scadenza"`
	Stato string    `gorm:"type:varchar(20);default:'Non usato'"`
}
type NotificheEmail struct{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TipoEvento string `gorm:"type:varchar(100);not null"`
	Destinatario string  `gorm:"not null"`
	Oggetto string  `gorm:"not null"`
	Corpo string `gorm:"type:text"`
}