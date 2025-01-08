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
        "private_key_id": "682d928919bc8404e0313cb1067163f05aec4fd7",
        "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCRnyOqXIZPFFKG\ngyUmVoYw3PEjxkY+crONALVrh+HJLqfrhkF/hB/2Xc362QjtrCgipMTsbG388kyX\ndEsbEVqGA4ObsR6xQ0R+k/t2c8ERZ/kZq1N8KgKS+vlxWSgevpICX6pbQpOdPWFL\n3/4wSKRIR6CPu5Jpm1Kb5X/8cCXNycERyuZkzwkVQ6G95PzVzkN7pRBeMNXZevuT\nncpcDAK4savUjfw6QoV8TcmMBE/1STR0ElypRG/9VTLCfn09lMtOkomBQfvxBkwr\nEQgxok5ZArl2qkbg25RnQSX96QFIC9a6wxat6ueCnr6EHvECebnXEnQsX5n4ODoQ\nb8l7aawLAgMBAAECggEADJid6EOzNtruRB/K+MGwpabIhKYtd4OXcwGiyVrj7z5f\nnKvGNaAGd1qbRL2j6Pi+MTgpQPDxVwxxQjaQ52LJ/q0GJAafPe99WMLQx4KKxEhg\nnG8cvEqr/CFee1OeRbTvtbanUNgo8Aezrb41FnsdTsbdvu3kMJaUmWyYJ6XE9100\nKvb/43NlBXW29hHhPayZrxPies/u9UuuTAObfynAeZuajYabNEKp5Jdga47nJoKO\n0tB5ZD3BGlNxzxKH1h4YaQCDH97Xmn4P1W4TekumaXuUCc77/6xmB5ampA49enpg\nVMvSyYpjVQL7wY4cfy99hbcj3L+iOnWkwsnfiMcjIQKBgQDKivNdAToREcF2nmPo\nbKGiCcMvrEBkkcv5uCBtXAKt7ZI/JK7ZCbZVYImdTYPsNdTirkdoiMnnaCNPiNmk\n6JOt3eKPUa8bLmW60i6wh4fZBq743QC6F+pky4lQsR/ZFwjBI2x7Ezn+i4IF9Fij\n3ZAa5ewCEWL07je/HGxSMWR1qwKBgQC4Dj4fCJVbwfoi0tdL0pYcEdhFFmO74qX6\n+bWuH+ApsPoSMlwFvRitvMbCiP+tkXNrVES8vTMiu5mlqNq6XNHkLdoAMVYgjxI1\nWIbmeby2ivO52qh9BsYExp6lmmH9j5OL6HtBPSaG6rKG+uyP23QbzjL1z6q7fy0l\nU4y4X8CDIQKBgDUxoRDADbwF6cV5e8vDHAAuiDCxEIhZMjT3gqy2CY4cYthqjfE6\nd8ScggfqH8edq7eNBfwSUNSRqRRuYJrK6l4zdBkn3tFIsjcKlHCQZ8E73CBICrTV\nKx4lxn6GxlKBli8DWq5IMmcDxLZDojQHMJ2f3Qf+APtKxSFQGbLMfhHZAoGBALYw\n2wtLKtXGgP2Bqb6TeSXWADf7PsRYSabTEiWHxhRe7Fug3/iKQ2iPakxc4oKEbTT8\nGIKf4oNqImCacdFyWg492QLB05itJwAJXpe8P7KOf04lBQ2l69QEbDxPQtqFCi++\n9GsSxhVdM1VsA0kvmZKAnW83nrC05hKBztUfa2ghAoGAXqNKYytt0+Ag2Zkm/BSn\nufkKgM1j6FdF/KGFTAcFQxfZ9mt4VPyg+8m/IhCTUX+G2jYzYAaOiV+z7iPPY5LH\niLo4xxQwPJJkmm8VzuPGh398rQv052SvyaaNY68wBNtz8oaUMGbpqjm+BjXuxYnd\nzdbVF2AKWY2jaN/Q1A3/Sv0=\n-----END PRIVATE KEY-----\n",
        "client_email": "firebase-adminsdk-yla8e@notification-push-40d24.iam.gserviceaccount.com",
        "client_id": "108746305878631738217",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-yla8e%40notification-push-40d24.iam.gserviceaccount.com",
        "universe_domain": "googleapis.com"
    }`

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
    return app
}