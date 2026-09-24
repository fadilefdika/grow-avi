package notification

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

type Payload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

func SendPushNotification(db *sql.DB, npk string, title, body, url string) {
	// Ambil semua subscription milik NPK ini
	rows, err := db.Query("SELECT endpoint, p256dh, auth FROM push_subscriptions WHERE npk = @p1", npk)
	if err != nil {
		log.Println("Error querying push subscriptions:", err)
		return
	}
	defer rows.Close()

	payloadData := Payload{
		Title: title,
		Body:  body,
		URL:   url,
	}
	payloadBytes, _ := json.Marshal(payloadData)

	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")

	if vapidPublicKey == "" || vapidPrivateKey == "" {
		log.Println("WARNING: VAPID keys not set, cannot send push notification")
		return
	}

	for rows.Next() {
		var endpoint, p256dh, auth string
		if err := rows.Scan(&endpoint, &p256dh, &auth); err != nil {
			continue
		}

		sub := &webpush.Subscription{
			Endpoint: endpoint,
			Keys: webpush.Keys{
				P256dh: p256dh,
				Auth:   auth,
			},
		}

		resp, err := webpush.SendNotification(payloadBytes, sub, &webpush.Options{
			Subscriber:      "mailto:admin@grow-avi.local", 
			VAPIDPublicKey:  vapidPublicKey,
			VAPIDPrivateKey: vapidPrivateKey,
			TTL:             3600,
		})

		if err != nil {
			log.Println("Failed to send push notification:", err)
		} else {
			resp.Body.Close()
			log.Printf("Push notification sent to NPK: %s", npk)
		}
	}
}
