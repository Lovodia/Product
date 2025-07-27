package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "host=localhost port=5432 user=appuser password=12345 dbname=products sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error opening connection:", err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Println("Error connection to BD", err)
		return err
	}

	log.Println("Connected to database successfully")
	return nil
}
