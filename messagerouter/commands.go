package messagerouter

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"strings"
	"sync"
)

var lock = new(sync.Mutex)
var commands = make(map[string]MessageHandler)

type CommandRequest struct {
	Username string
	Text     string
	Message  *slack.Msg
}

// MessageHandler responds to a slack or cli message request.
type MessageHandler func(*CommandRequest, ResponseWriter) error

// GetUsername retrieves human readable username information from SLACK API using Slack username ID
func (mr *MessageRouter) GetUsername(msg *slack.Msg) (string, error) {
	if mr.IsDev {
		return "cliuser", nil
	}
	userInfo, err := mr.API.GetUserInfo(msg.User)
	if err != nil {
		log.Println("No user information: "+msg.Text, err.Error(), msg)
		return "", err
	}

	return userInfo.Name, nil
}

// ProcessMessage performs basic routing on Slack messages in the channel #puddle
// Eventually, this will feed into two routers
// One for admin commands & one for mud commands
func (mr *MessageRouter) ProcessMessage(msg *slack.Msg) error {
	if hasPrefix(msg.Text) {
		username, _ := mr.GetUsername(msg)

		words := strings.Split(msg.Text, " ")
		cmdString := strings.ToLower(words[0])

		var w ResponseWriter
		if mr.IsDev {
			w = new(CLIResponseWriter)
		} else {
			w = new(SlackResponseWriter)
		}

		cr := CommandRequest{username, msg.Text, msg}

		if fn, ok := commands[cmdString]; ok {
			err := fn(&cr, w)
			if err != nil {
				w.WriteError(err.Error())
			}
			return err
		} else {
			return errors.New("Invalid command. Get it together")
		}
	}
	return nil
}

// hasPrefix checks for command prefix
func hasPrefix(s string) bool {
	return strings.HasPrefix(s, "!") || strings.HasPrefix(s, "%") || strings.HasPrefix(s, "#")
}

// Handle registers the handler for the given command
// Must match a command prefix, `!` for admin commands, `%` for general commands & `#` for MUD commands
func (mr *MessageRouter) Handle(cmdString string, fn MessageHandler) error {
	if !hasPrefix(cmdString) {
		log.Fatal(cmdString, " does not have a command prefix")
	}
	log.Println(fmt.Sprintf("Registering new admin command <%s>", cmdString))

	lock.Lock()
	defer lock.Unlock()

	if _, ok := commands[cmdString]; !ok {
		commands[cmdString] = fn
		return nil
	}

	return errors.New("Command exists")
}
