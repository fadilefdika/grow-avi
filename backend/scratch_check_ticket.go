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

	growID := "GR2609160001"
	var status string
	var isSpent int
	var deletedAt sql.NullString
	err = db.QueryRow("SELECT status, is_spent, CONVERT(varchar, deleted_at, 120) FROM activity_submissions WHERE grow_id = @p1", growID).Scan(&status, &isSpent, &deletedAt)
	if err != nil {
		log.Fatal("Error query:", err)
	}

	fmt.Printf("Ticket %s: Status=%s, IsSpent=%d, DeletedAt=%v\n", growID, status, isSpent, deletedAt.String)
}
