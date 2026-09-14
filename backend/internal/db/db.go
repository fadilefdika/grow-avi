package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	connString := os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		log.Fatal("DB_CONNECTION_STRING is not set in .env")
	}

	var err error
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: %s", err.Error())
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %s", err.Error())
	}

	fmt.Println("Connected to SQL Server successfully!")
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}
