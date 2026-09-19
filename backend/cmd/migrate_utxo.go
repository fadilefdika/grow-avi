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

	// Alter activity_submissions to add is_spent
	_, err := sqlDB.Exec("ALTER TABLE activity_submissions ADD is_spent BIT DEFAULT 0")
	if err != nil {
		fmt.Println("Error adding is_spent to activity_submissions:", err)
	} else {
		fmt.Println("Successfully added is_spent to activity_submissions")
	}

    // Since we also need to map a redemption to specific GROW IDs, we should create a mapping table:
    // redemption_items (redemption_id, submission_id)
    createTableQuery := `
    IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='redemption_items' and xtype='U')
    CREATE TABLE redemption_items (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        redemption_id BIGINT NOT NULL,
        submission_id BIGINT NOT NULL,
        FOREIGN KEY (redemption_id) REFERENCES reward_redemptions(id),
        FOREIGN KEY (submission_id) REFERENCES activity_submissions(id)
    )`
    _, err = sqlDB.Exec(createTableQuery)
    if err != nil {
        fmt.Println("Error creating redemption_items:", err)
    } else {
        fmt.Println("Successfully created redemption_items")
    }
}
