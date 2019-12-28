package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"github.com/peterhellberg/cloudmailin"
)

var (
	username     = os.Getenv("USERNAME")
	password     = os.Getenv("PASSWORD")
	slackToken   = os.Getenv("SLACK_TOKEN")
	slackChannel = os.Getenv("SLACK_CHANNEL")
)

func main() {
	handler := basicAuth(incomingMail, username, password, "Authentication required")
	http.HandleFunc("/incoming", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func basicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		handler(w, r)
	}
}

func incomingMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only POST requests are allowed using cloudmailin json format")
		return
	}
	// decode email content from request body
	defer r.Body.Close()
	msg, err := cloudmailin.Decode(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Failed to parse email: %s", err)
		w.WriteHeader(400)
	}
	// send message to slack channel
	channelID, timestamp, err := sendMessage(&msg)
	if err != nil {
		fmt.Fprintf(w, "Failed to send message to channel %s: %s", slackChannel, err)
		w.WriteHeader(500)
	} else {
		log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}
}

func sendMessage(data *cloudmailin.Data) (string, string, error) {
	api := slack.New(slackToken)
	attachment := slack.Attachment{
		AuthorName: data.Headers.From,
		Title:      data.Headers.Subject,
		Text:       data.Plain,
	}
	return api.PostMessage(slackChannel, slack.MsgOptionAttachments(attachment))
}
