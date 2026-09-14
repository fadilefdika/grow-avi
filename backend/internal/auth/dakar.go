package auth

import (
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

// ValidateCredentials checks credentials.
// For Awork, we fetch the users list and verify if the NPK exists.
// We expect password to match the default AWORK_DEFAULT_PASSWORD.
func ValidateCredentials(npk, password string) (bool, *DakarProfile, error) {
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

	// Assuming the response is {"data": [ { "npk": "...", "fullname": "...", "department": "..." } ]}
	type aworkUser struct {
		NPK        string `json:"npk"`
		Name       string `json:"fullname"`
		Department string `json:"department"`
	}
	var apiResp struct {
		Data []aworkUser `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return false, nil, err
	}

	for _, user := range apiResp.Data {
		if user.NPK == npk {
			fmt.Printf("[DEBUG DAKAR] Berhasil menemukan NPK: %s, Name: '%s', Dept: '%s'\n", user.NPK, user.Name, user.Department)
			return true, &DakarProfile{
				NPK:        user.NPK,
				UserName:   user.Name,
				Department: user.Department,
			}, nil
		}
	}

	fmt.Printf("[DEBUG DAKAR] NPK %s tidak ditemukan dalam daftar %d karyawan dari Awork\n", npk, len(apiResp.Data))
	return false, nil, nil
}
