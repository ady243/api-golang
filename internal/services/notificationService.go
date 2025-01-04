package services

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/ady243/teamup/internal"
)

type NotificationService struct {
	app    *firebase.App
	client *messaging.Client
}

func NewNotificationService() (*NotificationService, error) {

	app, err := internal.InitializeFirebase()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de Firebase: %v", err)
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("Erreur lors de la récupération du client de notifications: %v", err)
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

	// Envoie de la notification
	response, err := ns.client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("échec de l'envoi de la notification: %w", err)
	}

	log.Printf("Notification envoyée avec succès: %s", response)
	return nil
}
