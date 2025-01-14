package internal

import (
	"context"
	"log"
	"sync"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	app  *firebase.App
	once sync.Once
)

func InitializeFirebase() (*firebase.App, error) {
	credentials := `{
		  "type": "service_account",
		  "project_id": "notification-push-40d24",
		  "private_key_id": "76f0bbd328c6edbd345727fb790ff18b83b5b03f",
		  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDsdwzuKByQ93CJ\nO6mYXaATPcl2ep8L6jzf5PJjN689a72Vs1z1DQSklvS6YyHCADstEFGh9ewmUmzZ\nrXyjYG7nCp+PvMP/9n7W9jpXbHIUu4MJ1oe+4PNlM42WRudzb0588TY3dgRRvKq3\nDRk1xqQ3pPZZRJTkIMCxqxX4w3abNC+2vAhUWR/HUYoetz+3Lg5R9BoQfHmY7VH5\niNUWn8WlIncByXnMOMC+BzLotrKP995YWO33EPUk92PcbepCwdG1mijL3ynegQsO\n93JfHPWd7RepIqerLd4XMDHJIMeJAhLpDwyMbw81fBBS0fxnlaooTRXZO78ChLYT\nL/Npt0lFAgMBAAECggEABFO+8RswsUzfpws9gt6pgy7Zb4zgx6X7g4db9UKCxWBa\n2oLVsSPZH5/tHG9i90IIYWkBpk4jtPq1RHFfHug/IRLkhrZywqMQ1lwUI5DbaVp9\nvv4s2Ib4hF6F+UQxOtvSoeySE1yE5gx1Ep8Ey+pnfTwTZOCvrtIiZJNCaCvNu+JP\no1JVG/4GsfGRsrnTAUBsx2ecAOS/WkZTUiNIX20H6OTpWJqdKxc5o3FJZi/k8sI1\nKD6R9NuxZbGBuxbsiF8dTjn9D8ESJJGMPIv0ixbuhJxtJMfjHOE5rAYDQ/WNgA9s\nH/+IxBW4UNp7DggIQ8PXCBPg7MHuU10kz2UlJp1OMwKBgQD+jLTY8rks5fbBuinp\nLVgQU4ySZrOx07k+WCXvxYEM6fAlJMT8RucWlGdxo/eOU866eTez8LsCtew3hNF6\nsJnK3p+CIa65ZFG2kh9moZBxxXC9vDHYihCuPKGXixr9Chu1Q33SZNfxW0ClSQr2\n4vB9MCwB7S13UFL5z2Y31jq4hwKBgQDtz/chghbp6NQkIVeplSQBhETQhjOBkn3w\nErlqrJLI+GdELL4AbI5z0AiTJIW9JdIW4mba927CwtwY650iFJwBYgvF8ccC3ta/\n2cP/Lx2reCGonPgKxUDV+wv8NyLCYj3JS7+3ogaHgBmLcjyNVcGS7Tz1MKrG+ZeI\nCpb9evu+0wKBgFHiPYLUgdD5oOks07KYzY1i8wNdWkzICP0PKhT5ecwHrSKls2Bc\nBpZy4tvhnQ8B0qyVtd+CfwYeM4CgjypiiPaDqtgXsbcdmFOcqdFAA9E1bFD8qyQ3\nNap3ApxXOTVQ/RzQOzdlDTos2pzQ5GALHMWIq39rJocNJcQKfZ1UossdAoGBAK8A\nAekrlPcOecYYryzA7l0bW5RjnVV1Wq2m6cEhO2cevMdDcZJYUD/TT+wPzUbSpRZo\nBq6NtHkn8dV41Qn2RpMR9n30nLF1EGzfsEaCAoBjB8nPsQwj+cE9W6V/YVnP943A\n61UTq2BdGO8v4nVTLP6VC+2WoaWImETpHhFsRgM3AoGAIzJCcdeoXZ6Ql2Zvz5HR\nZLyJ63reKBzZoCIV7yjuzH0KzZsODCnVYRmUHvAYI+Q4jxs0fLI8iCTjJjH3t92g\noG5xF1WTTn3Vt/V4nXbIhNW/ctNEpySvLVztjjPKrF4bT1czy+Qy0vvKZzBsLm5X\nEs8dmMCr+S7SqA123nYhmao=\n-----END PRIVATE KEY-----\n",
		  "client_email": "firebase-adminsdk-yla8e@notification-push-40d24.iam.gserviceaccount.com",
		  "client_id": "108746305878631738217",
		  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
		  "token_uri": "https://oauth2.googleapis.com/token",
		  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-yla8e%40notification-push-40d24.iam.gserviceaccount.com",
		  "universe_domain": "googleapis.com"
		}
		`

	opt := option.WithCredentialsJSON([]byte(credentials))
	var err error
	app, err = firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "notification-push-40d24",
	}, opt)
	if err != nil {
		log.Fatalf("erreur lors de l'initialisation de l'application: %v", err)
	}
	return app, nil
}

func GetFirebaseApp() *firebase.App {
	InitializeFirebase()
	return app
}
