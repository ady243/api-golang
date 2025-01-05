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

	// Utilisation d'une variable d'environnement pour le chemin du fichier de credentials
	credentialsFilePath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	log.Printf("Valeur de FIREBASE_CREDENTIALS_PATH: %s", credentialsFilePath)
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
