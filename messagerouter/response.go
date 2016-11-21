package messagerouter

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
)

type ResponseWriter interface {
	Write(string) error
	WriteChannel(string, string) error
	WriteError(string) error
}

type SlackResponseWriter struct {
	rtm *slack.RTM
	msg *slack.Msg
}

func NewSlackResponseWriter(msg *slack.Msg, rtm *slack.RTM) *SlackResponseWriter {
	w := SlackResponseWriter{}
	w.msg = msg
	w.rtm = rtm

	return &w
}

//Msg sets the SlackResponseWriter Message
func (w *SlackResponseWriter) Msg(msg *slack.Msg) {
	w.msg = msg
}

//Msg sets the SlackResponseWriter RTM
func (w *SlackResponseWriter) Rtm(rtm *slack.RTM) {
	w.rtm = rtm
}

// WriteChannel sends message to particular channel
func (w *SlackResponseWriter) WriteChannel(channel string, text string) error {
	w.rtm.SendMessage(w.rtm.NewOutgoingMessage(text, channel))
	return nil
}

// Write writes to Slack
func (w *SlackResponseWriter) Write(text string) error {
	w.rtm.SendMessage(w.rtm.NewOutgoingMessage(text, w.msg.Channel))

	return nil
}

// WriteError writes an error to Slack
func (w *SlackResponseWriter) WriteError(errText string) error {
	lines := strings.Split(errText, "\n")
	for _, value := range lines {
		w.Write(":poop: " + value + " :poop:")
	}
	return nil
}

//CLIResponseWriter handles writing to the command line
type CLIResponseWriter struct {
}

// Write writes to CLI, prints channel it would go to if via Slack
func (w *CLIResponseWriter) WriteChannel(channel string, text string) error {
	fmt.Println(">> "+channel, text)

	return nil
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
