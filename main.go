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
	defer r.Body.Close()
	if msg, err := cloudmailin.Decode(r.Body); err == nil {
		sendMessage(&msg)
	} else {
		fmt.Fprintf(w, "Failed to parse request: %s", err)
		w.WriteHeader(400)
	}
}

func sendMessage(data *cloudmailin.Data) {
	api := slack.New(slackToken)
	attachment := slack.Attachment{
		AuthorName: data.Headers.From,
		Title:      data.Headers.Subject,
		Text:       data.Plain,
	}
	channelID, timestamp, err := api.PostMessage(slackChannel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
