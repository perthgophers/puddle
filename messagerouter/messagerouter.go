package messagerouter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
)

// MessageRouter does stuff with messages.
type MessageRouter struct {
	API     *slack.Client
	RTM     *slack.RTM
	token   string
	channel string
	gittag  string
	IsDev   bool
}

//Returns a new Pudddle Messagerouter
func New(token, gittag, channel string) *MessageRouter {
	mr := MessageRouter{
		token:   token,
		gittag:  gittag,
		channel: channel,
		IsDev:   false,
	}
	if os.Getenv("PUDDLEDEV") != "" || token == "" {
		mr.IsDev = true
	}

	if !mr.IsDev {
		mr.API = slack.New(token)

		// If you set debugging, it will log all requests to the console
		// Useful when encountering issues
		mr.API.SetDebug(true)

		mr.RTM = mr.API.NewRTM()
	}

	return &mr
}

func (mr *MessageRouter) Run() {

	if mr.IsDev {
		mr.runCLI()
	} else {
		mr.runSlack()
	}

	// RunCLI Starts the command line input shell

}

func (mr *MessageRouter) runSlack() {
	go mr.RTM.ManageConnection()
Loop:
	for {
		select {
		case msg := <-mr.RTM.IncomingEvents:
			log.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				log.Println("######### Connected to Slack #########")
				// log.Println("Infos:", ev.Info)
				// log.Println("Connection counter:", ev.ConnectionCount)
				// Replace #general with your Channel ID
				mr.RTM.SendMessage(mr.RTM.NewOutgoingMessage(fmt.Sprintf("... and I'm back! Git tag: %s", mr.gittag), mr.channel))

			case *slack.MessageEvent:
				j, _ := json.Marshal(ev.Msg)
				log.Printf("Message: %v\n", string(j))
				mr.ProcessMessage(&ev.Msg)

			case *slack.PresenceChangeEvent:
				log.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				log.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// log.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}

func (mr *MessageRouter) runCLI() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(fmt.Sprintf("... and I'm back! Git tag: %s", mr.gittag))
	fmt.Print("Puddle> ")
	for scanner.Scan() {
		line := scanner.Text()
		msg := slack.Msg{
			Text: line,
		}
		mr.ProcessMessage(&msg)
		fmt.Print("Puddle> ")
	}
}
