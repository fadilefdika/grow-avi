package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	connString := "server=localhost;database=grow_point;trusted_connection=yes;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()


	// Count all tickets by status for Veronica
	rows, err := db.Query("SELECT status, is_spent, COUNT(*) FROM activity_submissions WHERE npk = '0479' AND deleted_at IS NULL GROUP BY status, is_spent")
	if err != nil {
		log.Fatal("Error query:", err)
	}
	defer rows.Close()

	fmt.Println("All tickets for NPK 0479:")
	for rows.Next() {
		var status string
		var isSpent int
		var count int
		rows.Scan(&status, &isSpent, &count)
		fmt.Printf("Status: %s, IsSpent: %d, Count: %d\n", status, isSpent, count)
	}
}
