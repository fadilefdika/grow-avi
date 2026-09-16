package sync

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AworkUser struct {
	NPK        string `json:"npk"`
	Name       string `json:"fullname"`
	Department string `json:"department"`
}

// RunMassSync fetches all users from Awork API and upserts them into the local users table.
func RunMassSync(db *sql.DB) error {
	apiURL := os.Getenv("AWORK_API_URL")
	apiKey := os.Getenv("AWORK_API_KEY")

	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("AWORK_API_URL or AWORK_API_KEY not configured")
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Awork API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var apiResp struct {
		Data []AworkUser `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return err
	}

	if len(apiResp.Data) == 0 {
		return fmt.Errorf("no users returned from Awork API")
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Upsert logic for SQL Server
	// We use MERGE statement for efficient upsert
	stmt, err := tx.Prepare(`
		MERGE INTO users AS target
		USING (SELECT @p1 AS npk, @p2 AS fullname, @p3 AS department) AS source
		ON target.npk = source.npk
		WHEN MATCHED THEN
			UPDATE SET fullname = source.fullname, department = source.department, updated_at = GETDATE()
		WHEN NOT MATCHED THEN
			INSERT (npk, fullname, department) VALUES (source.npk, source.fullname, source.department);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, user := range apiResp.Data {
		_, err := stmt.Exec(user.NPK, user.Name, user.Department)
		if err != nil {
			return fmt.Errorf("failed to upsert user %s: %v", user.NPK, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("[SYNC] Berhasil melakukan mass sync %d users dari Awork API ke database lokal\n", len(apiResp.Data))
	return nil
}

// StartDailyCron runs the sync immediately on startup, and then once every 24 hours.
func StartDailyCron(db *sql.DB) {
	go func() {
		// Tunggu sebentar agar server nyala dulu
		time.Sleep(5 * time.Second)

		fmt.Println("[SYNC] Menjalankan sinkronisasi massal awal (Startup)...")
		if err := RunMassSync(db); err != nil {
			fmt.Printf("[SYNC ERROR] %v\n", err)
		}

		// Jadwal setiap 24 jam (sehari sekali)
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("[SYNC] Menjalankan sinkronisasi massal (Jadwal Harian)...")
			if err := RunMassSync(db); err != nil {
				fmt.Printf("[SYNC ERROR] %v\n", err)
			}
		}
	}()
}
