package messagerouter

import (
	"fmt"
	"github.com/nlopes/slack"
)

type ResponseWriter interface {
	Write(string) error
	WriteError(string) error
}

type SlackResponseWriter struct {
	rtm *slack.RTM
	msg *slack.Message
}

// SendMessage send slack or cli message
// shorthand for rtm.SendMessage(rtm.NewOutgoingMessage(text, CHANNEL))
func (w *SlackResponseWriter) Write(text string) error {
	w.rtm.SendMessage(w.rtm.NewOutgoingMessage(text, w.msg.Channel))

	return nil
}

// SendMessage send slack or cli a poopy error message
func (w *SlackResponseWriter) WriteError(errText string) error {
	return w.Write(":poop: " + errText + " :poop:")
}

//CLIResponseWriter handles writing to the command line
type CLIResponseWriter struct {
}

// SendMessage send slack or cli message
// shorthand for rtm.SendMessage(rtm.NewOutgoingMessage(text, CHANNEL))
func (w *CLIResponseWriter) Write(text string) error {
	fmt.Println(">> ", text)

	return nil
}

// SendMessage send slack or cli a poopy error message
func (w *CLIResponseWriter) WriteError(errText string) error {
	return w.Write(":poop: " + errText + " :poop:")
}
