package internal

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirebase() (*firebase.App, error) {
	ctx := context.Background()

	credentialsFilePath, err := filepath.Abs("./fireJson/notification-push-40d24-firebase-adminsdk-yla8e-95ee6971eb.json")
	if err != nil {
		log.Fatalf("Erreur lors de la résolution du chemin absolu: %v", err)
	}

	log.Printf("Chemin absolu du fichier de credentials: %s", credentialsFilePath)

	opt := option.WithCredentialsFile(credentialsFilePath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'initialisation de l'application: %v", err)
	}

	return app, nil
}
