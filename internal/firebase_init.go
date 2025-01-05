package internal

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var app *firebase.App

func InitializeFirebase() (*firebase.App, error) {
	credentials := `{
        "type": "service_account",
        "project_id": "notification-push-40d24",
        "private_key_id": "95ee6971eb02024b2f063536e9fdcb9b4aca0e61",
        "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9O++es7vtpQr6\n+AcIugkcilAKwcX8SFuQlWlUwc58ignpKP5t6oPmXJN/9nQBrlKn1TTvW75IrzOg\nf+jl9qo8REB2Ic3YYbzTjm1kCOAadhtxFRwbpyw5EDrVfAtJ91KM+T36D5xzgaUS\nvONYsRoa6lmObAfMKo4Z2D2nOVJU7DOcLiVyxqC0PclwHkN+sGRKa+663mxY2aVG\n+PCLvi69yxYriPs5qF95wt8zT16n8QIYS9tpIUhAfbATE6lOHEtkwZpCk/vqWCkI\nja6yNvsaaJz1bvQmiEocHnoV8AmBXQFJLEOmtw5Xfy7ucKof0pEBUREbfL95wBed\ntgebX7RjAgMBAAECggEAXF4krQUTXsD4Vp95o39XUjiLMGz8TJQvKahcpTQCq9Sf\nWNGbO6DEqE2Y69WrM2TZULXn5EwWWhk27Ily78kSuF8iTedbOFsg1e0IJVOVvCTZ\naT3CHhdgJSxwY2NsiTqxb0F7yJMVLWZjYn2TxWeRFAE/HJ9LwRMmkKP2GCmJMAzS\nEqgeh0vkmHDFI1WSziZWtfnCn7ddEHgkC+m8PzjbyJGI+jLKUO3fMPUEQAUOKjN/\nNn7D0LyBCHPsCDw4/98dpLr/fp8XHzNWgtzqBrh1R70xU9dzlIKARfnZ2s/qilY+\nyvYF5o63br/sx5roETqtdQipUk4kUdpC89cOQPX9LQKBgQDr7qo8/tKwuW8LpIrz\nxOLeP/Pn0QVvtfX5weAm0FDFGr4wvER0SLz9xPD/5FcFzILj7wxVDplnpCoKsuX+\nXCMGJbTBYS/sKx7R6PZCBWUQMBdeJckTuA09J+j9fLGDr7LAfn2Xs+P6vTBh+qrR\nph/Gf+sMAzUI3ldzTCco4JvgHwKBgQDNVG++cZ7pDwwb4JOeUFTAyiboV/fER3Ur\na+Pq2fHO29SY7oHd++CnNumvu4oOs+hrK6MZWIogwpJXup0wwqc8Rs24fitsy5B5\nqDmVJbQq+YjZ36zLq0er5fWSM0z05sQEOeiFDyMJ9LPemQiYGOBts4m7l0tKG4Wr\nQIZBJEcTPQKBgARZSdoF+Gw5fsqAJe+IWYYvN5e2Sptch0QrRq5weIypiYfscHaU\nQKeM0cRluRTqSB9bcKbAtiMq63t3ALZHjH24hDRsTi3UPaUw3hkpcEt0F3osyCAM\n7HGMIsdJXRxISMszia0aK8GbayDjNfLXVQ6bnQGDrZ6UOphtdutR+I2RAoGAS6QS\nclcLEpJfhJmL6CNxxX/zK17UwLMOYAuj7+2QHgNv41Lh9rQGg7NADWQKLPZr5acy\ns3ChmgXzwWvW4pKi5xqySIf6WV74f0jQxbgZEkfQ+WpkDrevdI0HlW9ep02n4mKu\n3O3Bm7ZQ2O0JdHadnuwoDjpjGw7ehEMF8lN594UCgYEAyL0Yo1JbTwnp2aMN23aL\nYwg/4vSzj97cDdlw0gy8iJtoAt/4HbAajQreBKESeXfcsxcHK/cz/uj+1yzGKv6c\n1JnlaHvqwBggBjh3MrQvXVsfZydtt8fXi6dQYuRpJ0mt4mbLheaCNtLdCsKFEnLU\nn3WvSNs3Qy0Dj7V6myLTf3c=\n-----END PRIVATE KEY-----\n",
        "client_email": "firebase-adminsdk-yla8e@notification-push-40d24.iam.gserviceaccount.com",
        "client_id": "108746305878631738217",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-yla8e%40notification-push-40d24.iam.gserviceaccount.com",
        "universe_domain": "googleapis.com"
    }`

	opt := option.WithCredentialsJSON([]byte(credentials))
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "notification-push-40d24",
	}, opt)
	if err != nil {
		log.Fatalf("erreur lors de l'initialisation de l'application: %v", err)
	}
	return app, nil
}

func GetFirebaseApp() *firebase.App {
	return app
}
