package internal

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirebase() (*firebase.App, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile("./fireJson/notification-push-40d24-firebase-adminsdk-yla8e-95ee6971eb.json")

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "notification-push-40d24",
	}, opt)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de Firebase: %v", err)
		return nil, err
	}

	return app, nil
}


