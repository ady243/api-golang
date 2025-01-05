package internal

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var app *firebase.App

func init() {
	opt := option.WithCredentialsFile("./fireJson/notification-push-40d24-firebase-adminsdk-yla8e-95ee6971eb.json")
	var err error
	app, err = firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "notification-push-40d24",
	}, opt)
	if err != nil {
		log.Fatalf("erreur lors de l'initialisation de l'application: %v", err)
	}
}

func GetFirebaseApp() *firebase.App {
	return app
}
