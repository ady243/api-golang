package services

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type NotificationService struct {
	app    *firebase.App
	client *messaging.Client
}

func NewNotificationService() (*NotificationService, error) {
	opt := option.WithCredentialsJSON([]byte(`
		{
			"type": "service_account",
			"project_id": "notification-push-40d24",
			"private_key_id": "7cf0d9b6aa6750c059c3beabd57d16d3ccc2fde4",
			"private_key": "-----BEGIN PRIVATE KEY-----\nMIIE...Q0xq\n-----END PRIVATE KEY-----\n",
			"client_email": "firebase-adminsdk-yla8e@notification-push-40d24.iam.gserviceaccount.com",
			"client_id": "108746305878631738217",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-yla8e%40notification-push-40d24.iam.gserviceaccount.com"
		}
	`))

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v", err)
		return nil, err
	}

	return &NotificationService{
		app:    app,
		client: client,
	}, nil
}

func (ns *NotificationService) SendPushNotification(token, title, body string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := ns.client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	log.Printf("Successfully sent message: %s", response)
	return nil
}
