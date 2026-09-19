package main

import (
	"fmt"
	"grow-point/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	db.InitDB()
	sqlDB := db.GetDB()

	rows, err := sqlDB.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reward_redemptions'")
	if err != nil {
		fmt.Println("Error querying redemptions schema:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var colName, dataType string
		rows.Scan(&colName, &dataType)
		fmt.Printf("Column: %s, Type: %s\n", colName, dataType)
	}

    rows2, _ := sqlDB.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'activity_submissions'")
    defer rows2.Close()
	fmt.Println("---")
	for rows2.Next() {
		var colName, dataType string
		rows2.Scan(&colName, &dataType)
		if colName == "id" {
            fmt.Printf("Column: %s, Type: %s\n", colName, dataType)
        }
	}
}
