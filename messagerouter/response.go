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

// Write writes to Slack
func (w *SlackResponseWriter) Write(text string) error {
	w.rtm.SendMessage(w.rtm.NewOutgoingMessage(text, w.msg.Channel))

	return nil
}

// WriteError writes an error to Slack
func (w *SlackResponseWriter) WriteError(errText string) error {
	return w.Write(":poop: " + errText + " :poop:")
}

//CLIResponseWriter handles writing to the command line
type CLIResponseWriter struct {
}

// Write writes to CLI
func (w *CLIResponseWriter) Write(text string) error {
	fmt.Println(">> ", text)

	return nil
}

// WriteError writes error to CLI
func (w *CLIResponseWriter) WriteError(errText string) error {
	return w.Write(":poop: " + errText + " :poop:")
}
