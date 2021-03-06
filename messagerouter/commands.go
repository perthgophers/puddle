package messagerouter

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/nlopes/slack"
)

var lock = new(sync.Mutex)
var commands = make(map[string]MessageHandler)
var msgprocessing = make([]MessageHandler, 0)

type CommandRequest struct {
	Username string
	Text     string
	Message  *slack.Msg
}

// MessageHandler responds to a slack or cli message request.
type MessageHandler func(*CommandRequest, ResponseWriter) error

// usernameCache stores usernames to protect against slack API exhaustion
var usernameCache = map[string]string{"0": "nii236"}

// GetUsername retrieves human readable username information from SLACK API using Slack username ID
func (mr *MessageRouter) GetUsername(msg *slack.Msg) (string, error) {
	if username, ok := usernameCache[msg.User]; ok {
		return username, nil
	}

	userInfo, err := mr.API.GetUserInfo(msg.User)
	if err != nil {
		log.Println("No user information: "+msg.Text, err.Error(), msg)
		return "", err
	}

	return userInfo.Name, nil
}

// BuildResponseWriter checks the MessageRouter and returns the appropriate response writer
func (mr *MessageRouter) BuildResponseWriter(msg *slack.Msg) ResponseWriter {
	var w ResponseWriter
	if mr.IsDev {
		w = new(CLIResponseWriter)
	} else {
		w = NewSlackResponseWriter(msg, mr.RTM)
	}

	return w
}

// ProcessMessage performs basic routing on Slack messages in the channel #puddle
// Eventually, this will feed into two routers
// One for admin commands & one for mud commands
func (mr *MessageRouter) ProcessMessage(msg *slack.Msg) error {
	username, _ := mr.GetUsername(msg)
	cr := CommandRequest{username, msg.Text, msg}
	w := mr.BuildResponseWriter(msg)

	if hasPrefix(msg.Text) {
		words := strings.Split(msg.Text, " ")
		cmdString := strings.ToLower(words[0])

		if fn, ok := commands[cmdString]; ok {
			err := fn(&cr, w)
			if err != nil {
				w.WriteError(err.Error())
			}
			return err
		} else {
			return errors.New("Invalid command. Get it together")
		}
	} else {
		for _, handler := range msgprocessing {
			handler(&cr, w)
		}
	}
	return nil
}

// hasPrefix checks for command prefix
func hasPrefix(s string) bool {
	return strings.HasPrefix(s, "!") || strings.HasPrefix(s, "%") || strings.HasPrefix(s, "#")
}

// Handle registers the handler for the given command, or register a standard message parsing function
// Must match a command prefix, `!` for admin commands, `%` for general commands & `#` for MUD commands
// For functions that will parse all messages (sans command prefix), the cmdString must equal "*"
func (mr *MessageRouter) Handle(cmdString string, fn MessageHandler) error {
	if !hasPrefix(cmdString) && cmdString != "*" {
		log.Fatal(cmdString, " does not have a command prefix")
	}
	log.Println(fmt.Sprintf("Registering new admin command <%s>", cmdString))

	lock.Lock()
	defer lock.Unlock()

	if cmdString == "*" {
		msgprocessing = append(msgprocessing, fn)
		return nil
	}

	if _, ok := commands[cmdString]; !ok {
		commands[cmdString] = fn
		return nil
	}

	return errors.New("Command exists")
}
