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
	// Spécifie le chemin vers le fichier JSON téléchargé
	opt := option.WithCredentialsFile("./json/teamup-6b3b6-firebase-adminsdk-1zv3d-7b3b6b7b7e.json")

	// Initialise l'application Firebase
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
		return nil, err
	}

	// Crée un client pour envoyer des notifications push
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

	// Envoie de la notification
	response, err := ns.client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	log.Printf("Successfully sent message: %s", response)
	return nil
}
