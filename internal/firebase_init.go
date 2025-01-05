package internal

import (
    "context"
    firebase "firebase.google.com/go"
    "google.golang.org/api/option"
)

var app *firebase.App

func InitializeFirebase() (*firebase.App, error) {
    opt := option.WithCredentialsFile("./notification-push-40d24-firebase-adminsdk-yla8e-95ee6971eb.json")
    var err error
    app, err = firebase.NewApp(context.Background(), &firebase.Config{
        ProjectID: "notification-push-40d24",
    }, opt)
    if err != nil {
        return nil, err
    }
    return app, nil
}

func GetFirebaseApp() *firebase.App {
    return app
}