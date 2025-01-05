package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var app *firebase.App

func InitializeFirebase() (*firebase.App, error) {
	credentialsFilePath := "./notification-push-40d24-firebase-adminsdk-yla8e-95ee6971eb.json"

	absPath, err := filepath.Abs(credentialsFilePath)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du chemin absolu: %v", err)
	}
	log.Printf("Chemin absolu du fichier de clés de service Firebase: %s", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("le fichier de clés de service Firebase n'existe pas: %v", err)
	}

	opt := option.WithCredentialsFile(absPath)
	app, err = firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "notification-push-40d24",
	}, opt)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'initialisation de l'application: %v", err)
	}
	return app, nil
}

func GetFirebaseApp() *firebase.App {
	return app
}
