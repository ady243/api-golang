package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type NotificationService struct {
	OneSignalAppID  string
	OneSignalAPIKey string
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		OneSignalAppID:  os.Getenv("ONESIGNAL_APP_ID"),
		OneSignalAPIKey: os.Getenv("ONESIGNAL_API_KEY"),
	}
}

type OneSignalNotification struct {
	AppID            string            `json:"app_id"`
	Headings         map[string]string `json:"headings"`
	Contents         map[string]string `json:"contents"`
	IncludePlayerIDs []string          `json:"include_player_ids"`
}

func (s *NotificationService) SendNotification(playerID, title, message string) error {
	notification := OneSignalNotification{
		AppID:            s.OneSignalAppID,
		Headings:         map[string]string{"en": title},
		Contents:         map[string]string{"en": message},
		IncludePlayerIDs: []string{playerID},
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", s.OneSignalAPIKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification: %s", resp.Status)
	}

	return nil
}
