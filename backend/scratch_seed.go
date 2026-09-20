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

	growID := "GR999999999"

	// Delete from activity_submissions
	query := `DELETE FROM activity_submissions WHERE grow_id = @p1`
	res, err := db.Exec(query, growID)
	if err != nil {
		log.Fatal("Error deleting submission:", err)
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Printf("Successfully deleted %d dummy submission(s) with GROW ID %s\n", rowsAffected, growID)
}
