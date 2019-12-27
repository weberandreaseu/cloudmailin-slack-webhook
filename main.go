package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"github.com/peterhellberg/cloudmailin"
)

var (
	slackToken   = os.Getenv("SLACK_TOKEN")
	slackChannel = os.Getenv("SLACK_CHANNEL")
)

func main() {
	http.HandleFunc("/incoming", incomingMail)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
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
		fmt.Printf("Failed to parse request: %s", err)
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
