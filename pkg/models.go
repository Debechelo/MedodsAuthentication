package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitializeDB инициализирует подключение к базе данных
func InitializeDB() {
	var err error

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connection initialized")

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL,
		token TEXT NOT NULL,
		ip_address TEXT NOT NULL,
	    UNIQUE (user_id)
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error creating refresh_tokens table: %v", err)
	}
}

func SaveRefreshTokenToDB(userID, ip, hashedToken string) error {
	insertToken := `
	INSERT INTO refresh_tokens (user_id, token, ip_address)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id)
	DO UPDATE SET token = $2, ip_address = $3`

	_, err := db.Exec(insertToken, userID, hashedToken, ip)
	if err != nil {
		log.Printf("Error saving refresh token to DB for user %s: %v", userID, err)
		return err
	}

	log.Printf("Refresh token saved/updated for user %s", userID)
	return err
}

func GetRefreshTokenFromDB(userID string) (string, string, error) {
	var hashedToken, ipAddress string

	query := `
	SELECT token, ip_address 
	FROM refresh_tokens 
	WHERE user_id = $1`

	err := db.QueryRow(query, userID).Scan(&hashedToken, &ipAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No refresh token found for user %s", userID)
			return "", "", nil
		}
		log.Printf("Error fetching refresh token for user %s: %v", userID, err)
		return "", "", err
	}

	return hashedToken, ipAddress, nil
}
