package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	connString := os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		log.Fatal("DB_CONNECTION_STRING is not set")
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}

	fmt.Println("Connected to Database!")

	// 1. Tambah reward_type
	_, err = db.Exec("ALTER TABLE rewards ADD reward_type VARCHAR(50) DEFAULT 'CUSTOM'")
	if err != nil {
		fmt.Println("Error adding reward_type to rewards:", err)
	} else {
		fmt.Println("Successfully added reward_type to rewards")
	}

	// 2. Tambah is_locked
	_, err = db.Exec("ALTER TABLE activity_submissions ADD is_locked BIT DEFAULT 0")
	if err != nil {
		fmt.Println("Error adding is_locked to activity_submissions:", err)
	} else {
		fmt.Println("Successfully added is_locked to activity_submissions")
	}
}
