package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	connString := "server=localhost;database=grow_point;trusted_connection=yes;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: %s", err.Error())
	}
	defer db.Close()

	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	// 1. Update activities
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM activities WHERE default_points = 0 OR default_points IS NULL")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	
	type activity struct {
		id   int
		name string
	}
	var acts []activity
	for rows.Next() {
		var a activity
		rows.Scan(&a.id, &a.name)
		acts = append(acts, a)
	}
	rows.Close()

	for _, a := range acts {
		newPoints := (rand.Intn(10) + 1) * 5
		_, err := db.ExecContext(ctx, "UPDATE activities SET default_points = @p1 WHERE id = @p2", sql.Named("p1", newPoints), sql.Named("p2", a.id))
		if err != nil {
			fmt.Printf("Failed to update %s: %s\n", a.name, err.Error())
		} else {
			fmt.Printf("Updated activity '%s' to %d points.\n", a.name, newPoints)
		}
	}

	fmt.Println("Database update complete.")
}
