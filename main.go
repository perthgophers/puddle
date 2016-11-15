package main

import (
	"flag"
	"log"

	"github.com/nlopes/slack"
)

// SLACKTOKEN is the slack API token
var SLACKTOKEN string

func init() {
	flag.StringVar(&SLACKTOKEN, "slack_token", "", "Slack API Token")
}

func main() {
	flag.Parse()
	if SLACKTOKEN == "" {
		log.Fatalln("No slack token provided")
	}
	api := slack.New(SLACKTOKEN)

	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			log.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev)
			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop
			default:
			}
		}
	}
}
