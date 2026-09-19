package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	connString := "server=localhost;database=grow_point;trusted_connection=yes;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer db.Close()

	ctx := context.Background()
	
	rows, err := db.QueryContext(ctx, "SELECT id, name, default_points FROM activities")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var pts sql.NullInt64
		rows.Scan(&id, &name, &pts)
		if !pts.Valid || pts.Int64 == 0 {
			fmt.Printf("Activity %s (ID %d) has 0 or NULL points\n", name, id)
		}
	}
}
