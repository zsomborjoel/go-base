package email

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

const EmailTitle = "Workoutxz account confirmation"

func SendEmail(target string) error {
	account := os.Getenv("GMAIL_ACCOUNT")
	password := os.Getenv("GMAIL_PASSWORD")
	htmlPath := os.Getenv("EMAIL_HTML_PATH")

	html, err := os.ReadFile(htmlPath)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", account)
	msg.SetHeader("To", target)
	msg.SetHeader("Subject", EmailTitle)
	msg.SetBody("text/html", string(html))

	n := gomail.NewDialer("smtp.gmail.com", 587, account, password)
	if err := n.DialAndSend(msg); err != nil {
		log.Error().Err(err)
		return err
	}

	return nil
}
