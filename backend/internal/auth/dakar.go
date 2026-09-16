package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DakarProfile represents the profile data returned by Dakar API
type DakarProfile struct {
	NPK        string
	UserName   string
	Department string
}

type aworkUser struct {
	NPK        string `json:"npk"`
	Name       string `json:"fullname"`
	Department string `json:"department"`
}

// ValidateCredentials checks credentials.
func ValidateCredentials(db *sql.DB, npk, password string) (bool, *DakarProfile, error) {
	if npk == "" || password == "" {
		return false, nil, errors.New("empty credentials")
	}

	defaultPassword := os.Getenv("AWORK_DEFAULT_PASSWORD")
	if defaultPassword == "" {
		defaultPassword = "password" // fallback
	}

	// Check admin override
	if npk == "admin" && password == "admin" {
		return true, &DakarProfile{
			NPK:        "admin",
			UserName:   "Super Admin",
			Department: "Management",
		}, nil
	}

	if password != defaultPassword {
		return false, nil, nil // Invalid password
	}

	// 1. Cek di tabel users lokal (Fast Path)
	var fullname, department string
	err := db.QueryRow("SELECT fullname, department FROM users WHERE npk = @p1", npk).Scan(&fullname, &department)
	if err == nil {
		fmt.Printf("[DEBUG DAKAR] NPK %s ditemukan di database lokal\n", npk)
		return true, &DakarProfile{
			NPK:        npk,
			UserName:   fullname,
			Department: department,
		}, nil
	}

	if err != sql.ErrNoRows {
		// Error database selain tidak ditemukan
		return false, nil, fmt.Errorf("database error: %v", err)
	}

	// 2. JIT Fallback (Jika tidak ada di DB, tembak API Awork sekali ini saja)
	fmt.Printf("[DEBUG DAKAR] NPK %s tidak ada di DB lokal, mencoba JIT Fallback ke API Awork...\n", npk)
	apiURL := os.Getenv("AWORK_API_URL")
	apiKey := os.Getenv("AWORK_API_KEY")

	if apiURL == "" || apiKey == "" {
		// Mock fallback if env is missing
		if npk == "12345" && password == "password" {
			return true, &DakarProfile{
				NPK:        "12345",
				UserName:   "John Doe",
				Department: "IT",
			}, nil
		}
		return false, nil, errors.New("AWORK_API_URL or AWORK_API_KEY not configured")
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, nil, fmt.Errorf("Awork API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}

	var apiResp struct {
		Data []aworkUser `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return false, nil, err
	}

	// Cari NPK dari response Awork
	for _, u := range apiResp.Data {
		if u.NPK == npk {
			// Masukkan ke DB lokal (Upsert) agar next login langsung kena Fast Path
			_, errIns := db.Exec(`
				MERGE INTO users AS target
				USING (SELECT @p1 AS npk, @p2 AS fullname, @p3 AS department) AS source
				ON target.npk = source.npk
				WHEN NOT MATCHED THEN
					INSERT (npk, fullname, department) VALUES (source.npk, source.fullname, source.department);
			`, u.NPK, u.Name, u.Department)

			if errIns != nil {
				fmt.Printf("[DEBUG DAKAR] Gagal insert JIT user %s ke lokal: %v\n", npk, errIns)
			} else {
				fmt.Printf("[DEBUG DAKAR] Berhasil JIT sync untuk user %s\n", npk)
			}

			return true, &DakarProfile{
				NPK:        u.NPK,
				UserName:   u.Name,
				Department: u.Department,
			}, nil
		}
	}

	fmt.Printf("[DEBUG DAKAR] JIT Fallback gagal, NPK %s memang tidak ada di Awork\n", npk)
	return false, nil, nil
}
