package services

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/ady243/teamup/internal"
)

type NotificationService struct {
	app          *firebase.App
	client       *messaging.Client
	redisService *RedisService
}

func NewNotificationService(redisService *RedisService) (*NotificationService, error) {
	app := internal.GetFirebaseApp()
	if app == nil {
		return nil, fmt.Errorf("firebase app non initialisée")
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("Erreur lors de la récupération du client de notifications: %v", err)
		return nil, err
	}

	return &NotificationService{
		app:          app,
		client:       client,
		redisService: redisService,
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
		return fmt.Errorf("échec de l'envoi de la notification: %w", err)
	}

	log.Printf("Notification envoyée avec succès: %s", response)

	notificationKey := fmt.Sprintf("notification:%s", token)
	notificationValue := fmt.Sprintf("%s: %s", title, body)
	if err := ns.redisService.SetNotification(notificationKey, notificationValue, 14*24*time.Hour); err != nil {
		return fmt.Errorf("échec du stockage de la notification dans Redis: %w", err)
	}

	return nil
}

func (ns *NotificationService) GetUnreadNotifications(token string) ([]string, error) {
	keys, err := ns.redisService.GetKeys(fmt.Sprintf("notification:%s*", token))
	if err != nil {
		return nil, fmt.Errorf("échec de la récupération des notifications non lues: %w", err)
	}

	notifications := make([]string, 0, len(keys))
	for _, key := range keys {
		notification, err := ns.redisService.GetNotification(key)
		if err != nil {
			return nil, fmt.Errorf("échec de la récupération de la notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (ns *NotificationService) MarkNotificationsAsRead(token string) error {
	keys, err := ns.redisService.GetKeys(fmt.Sprintf("notification:%s*", token))
	if err != nil {
		return fmt.Errorf("échec de la récupération des notifications non lues: %w", err)
	}

	for _, key := range keys {
		if err := ns.redisService.DeleteNotification(key); err != nil {
			return fmt.Errorf("échec de la suppression de la notification: %w", err)
		}
	}

	return nil
}
