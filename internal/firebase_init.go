package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirebase() (*firebase.App, error) {
	ctx := context.Background()

	credentialsFilePath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credentialsFilePath == "" {
		log.Fatalf("Le chemin des credentials Firebase n'est pas défini")
	}

	opt := option.WithCredentialsFile(credentialsFilePath)

	// Configuration explicite de l'ID du projet
	config := &firebase.Config{
		ProjectID: "notification-push-40d24",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'initialisation de l'application: %v", err)
	}

	return app, nil
}
