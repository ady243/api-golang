

package internal

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirebase() (*firebase.App, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile("./services/json/teamup-6b3b6-firebase-adminsdk-1zv3d-7b3b6b7b7e.json") 

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de Firebase: %v", err)
		return nil, err
	}

	return app, nil
}
