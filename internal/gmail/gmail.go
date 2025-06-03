package gmail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Attachment represents a file attachment with its filename and data.
type Attachment struct {
	Filename string
	Data     []byte
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
	fmt.Printf("Go to the following link in your browser then type the " +
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
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

// Saves a token to a file.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetService() *gmail.Service {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return srv
}

// GetMessages fetches a list of Gmail messages filtered by a query.
func GetMessages(srv *gmail.Service, query string) ([]*gmail.Message, error) {
	var messages []*gmail.Message
	msgReq := srv.Users.Messages.List("me").Q(query)
	err := msgReq.Pages(context.Background(), func(resp *gmail.ListMessagesResponse) error {
		for _, m := range resp.Messages {
			msg, err := srv.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				return fmt.Errorf("unable to retrieve message: %v", err)
			}
			messages = append(messages, msg)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %v", err)
	}
	return messages, nil
}

// GetAttachments extracts PDF attachments from a Gmail message.
func GetAttachments(srv *gmail.Service, message *gmail.Message) ([]Attachment, error) {
	var attachments []Attachment

	if message.Payload.Parts == nil {
		return attachments, nil
	}

	for _, part := range message.Payload.Parts {
		if part.MimeType == "application/pdf" && part.Filename != "" {
			att, err := srv.Users.Messages.Attachments.Get("me", message.Id, part.Body.AttachmentId).Do()
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve attachment: %v", err)
			}

			data, err := base64.URLEncoding.DecodeString(att.Data)
			if err != nil {
				return nil, fmt.Errorf("unable to decode attachment data: %v", err)
			}

			attachments = append(attachments, Attachment{
				Filename: part.Filename,
				Data:     data,
			})
		}
	}

	return attachments, nil
}
