package email

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

const EmailTitle = "Workoutxz account confirmation"

func SendEmail(target string) error {
	log.Debug().Msg("email.SendEmail called")

	account := os.Getenv("GMAIL_ACCOUNT")
	htmlPath := os.Getenv("EMAIL_HTML_PATH")
	secretPath := os.Getenv("SECRET_PATH")

	log.Info().Msg(account)

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

	n := gomail.NewDialer("smtp.gmail.com", 587, account, "")
	if err := n.DialAndSend(msg); err != nil {
		log.Error().Err(err)
	}

	ctx := context.Background()
	b, err := os.ReadFile(secretPath)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read client secret file")
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse client secret file to config")
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve Gmail client")
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve labels")
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return nil
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}

	log.Info().Msg(fmt.Sprintf("Email sent to %s", target))

	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Info().Msg(fmt.Sprintf("Go to the following link in your browser then type the "+
		"authorization code: %s", authURL))

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Error().Err(err).Msg("Unable to read authorization code")
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve token from web")
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Info().Msg(fmt.Sprintf("Saving credential file to: %s\n", path))
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Error().Err(err).Msg("Unable to cache oauth token")
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
