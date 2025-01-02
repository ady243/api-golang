package services

import (
	"bytes"
	"net/smtp"
	"os"
	"text/template"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: white;
        }
        .container {
            font-family: Arial, sans-serif;
            margin: 0 auto;
            padding: 20px;
            background-color: white;
            border-radius: 5px;
            width: 100%;
            height: 100%;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            margin-top: 20px;
            font-size: 16px;
            color: #fff;
            background-color: #01BF6B;
            text-decoration: none;
            border-radius: 5px;
            cursor: pointer;
            width: 300px;
            height: 40px;
            justify-content: center;
            align-items: center;
            text-align: center;
        }

        .button-container {
            display: flex;
            justify-content: center;
            align-items: center;
            text-align: center;
        }

        .button:hover {
            background-color: #018966;
        }

        .button a {
            text-decoration: none;
            font-weight: bold;
            font-size: 38px;
        }

        p{
            font-size: 18px;
        }
        
        .button-text {
            font-weight: bold;
            font-size: 18px;
        }
        
        h1 {
            text-decoration: none;
            font-size: 20px;
        }
        
        h2 {
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Bonjour {{.ToEmail}},</h2>
        <h3>Bienvenue sur TeamUp 💫😁️!</h3>
        <p>Pour continuer l'aventure avec TeamUp⚽️, veuillez confirmer votre compte en cliquant sur le bouton ci-dessous ⌛.</p>
        <div class="button-container">
           <a href="http://localhost:3003/api/confirm_email?token={{.Token}}" class="button">
                <span class="button-text">Confirmer mon compte</span>
            </a>
        </div>
    </div>
</body>
</html>
`

type EmailData struct {
	ToEmail string
	Subject string
	Token   string
}

// SendConfirmationEmail envoie un email de confirmation à un utilisateur avec un lien de confirmation
//
// Il enverra un email à l'adresse email donnée avec un lien pour confirmer leur compte.
// Le lien contiendra un jeton qui peut être utilisé pour confirmer le compte.
//
// Paramètres:
// - toEmail: l'adresse email de l'utilisateur à qui envoyer l'email de confirmation
// - token: le jeton de confirmation à inclure dans le lien
//
// Retourne:
// - error: une erreur si l'email n'a pas pu être envoyé, nil sinon
func (e *EmailService) SendConfirmationEmail(toEmail, token string) error {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")
	host := "smtp.gmail.com"
	port := "587"

	auth := smtp.PlainAuth("", from, password, host)
	subject := "Confirmez votre compte"

	data := EmailData{
		ToEmail: toEmail,
		Subject: subject,
		Token:   token,
	}

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body.String())

	err = smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}
