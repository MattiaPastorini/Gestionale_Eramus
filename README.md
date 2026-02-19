# Eramus

# Gestionale Eramus

### Descrizione

Applicazione full-stack per la gestione completa di magazzino e utenti aziendali. Include dashboard interattiva con grafici, CRUD completo per prodotti e utenti, autenticazione JWT sicura e gestione automatica refresh token.

### Caratteristiche Principali

- Autenticazione JWT: con refresh token automatico
- Dashboard: interattiva con statistiche e grafici (Chart.js)
- Gestione Utenti: completa (CRUD + ruoli)
- Gestione Inventario: con soglia minima magazzino
- Responsive Design Bootstrap Italia
- Sicurezza: bcrypt + validazione password robusta

### Linguaggi e tecnologie usate

- Frontend: JavaScript con framework Next.js (App Router preferibile)
- Backend: Go con implementazione API REST
- Database: PostgreSQL con utilizzo di migrazioni
- Stile grafico: Bootstrap Italia
- Autenticazione: JWT (Access + Refresh Token)

### Requisiti

- Go 1.23+
- Node.js 20+ e npm
- PostgreSQL 17

### Backend .env

- Nel file `.env` vanno definiti:
- DB_NAME=
- DB_PASSWORD=
- JWT_SECRET=

### Struttura DB

# Utente

- Id (UUID - PK)
- Username (univoco)
- Email (univoca)
- Password (hash cifrato)
- Nome
- Cognome
- Data nascita
- Ruolo (FK)
- Tentativi login falliti
- Stato account (Attivo / Bloccato)
- Ultimo login
- Data creazione
- Data aggiornamento

# Ruolo

- Id (UUID - PK)
- Nome ruolo (Admin / Operatore)
- Descrizione

# Tipo Prodotto

- ld (UUID - PK)
- Corpo messaggio
- Data invio
- Esito invio

# Prodotto

- Nome oggetto
- Descrizione
- Quantità disponibile
- Prezzo unitario
- Soglia minima di magazzino
- Data inserimento
- Tipo prodotto (FK)
- Creato da (FK utente)
- Data ultima modifica

# Movimento Magazzino

- Id (UUID - PK)
- Prodotto (FK)
- Tipo movimento (Carico / Scarico)
- Quantità
- Data movimento
- Utente operazione (FK)
- Note

# Log Accessi

- Id (UUID - PK)
- Utente (FK)
- Data accesso
- Esito (Successo / Fallito)
- Indirizzo IP

# Recupero Password

- Id (UUID - PK)
- Utente (FK)
- Token univoco
- Data generazione
- Data scadenza (1 ora)
- Stato (Usato / Non usato)

# Notifiche Email

- Id (UUID - PK)
- Tipo evento (Nuovo Utente / Soglia Minima / Reset Password)
- Destinatario
- Oggetto
- Invio email automatica all'amministratore se quantità sotto soglia minima

### API Endpoints

# Metodo # Endpoint # Descrizione # Autorizzazione

- POST /api/login (Autenticazione)
- POST /api/refresh (Refresh token)
- GET /api/dashboard/statistiche (Dashboard stats)
- GET /api/utenti/ruoli (Lista ruoli)
- GET /api/utenti (Lista utenti)
- POST /api/utenti (Crea utente)
- PUT /api/utenti/:id (Aggiorna utente)
- DELETE /api/utenti/:id (Disattiva utente)
- GET /api/inventario/tipi (Tipi prodotti)
- GET /api/inventario/prodotti (Lista prodotti)
- POST /api/inventario/prodotti (Crea prodotto)
- PUT /api/inventario/prodotti/:id/stock (Aggiorna stock)
- DELETE /api/inventario/prodotti/:id (Elimina prodotto)

### Clone Repository

- git clone https://github.com/MattiaPastorini/Gestionale_Eramus.git

### Utilizzo

- Login: http://localhost:3000
- Dashboard: /dashboard
- Utenti: /utenti
- Inventario: /inventario

### Credenziali demo

- Admin: admin / Admin123!
- Operatore: user / User123! (non funzionante)
