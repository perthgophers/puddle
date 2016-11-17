package responses

import (
	"github.com/nlopes/slack"
	"github.com/perthgophers/puddle/messagerouter"
	"testing"
)

type ResponseTester struct {
	cr messagerouter.CommandRequest
	w  messagerouter.ResponseWriter
}
type CommandRequest struct {
	Username string
	Text     string
	Message  *slack.Msg
}
type TestResponseWriter struct {
	ExpectedMessage string
	ActualMessage   string
}

func (w *TestResponseWriter) Write(text string) error {
	w.ActualMessage = text

	return nil
}

func (w *TestResponseWriter) WriteChannel(errText string) error {
	return nil
}

func (w *TestResponseWriter) WriteError(errText string) error {
	return w.Write(":poop: " + errText + " :poop:")
}

func TestPing(t *testing.T) {

	cr := messagerouter.CommandRequest{
		Username: "fakeuser",
		Text:     "!ping",
		Message:  new(slack.Msg),
	}

	w := TestResponseWriter{
		ExpectedMessage: "pong!",
	}

	Ping(&cr, &w)

	if w.ExpectedMessage != w.ActualMessage {
		t.Error(
			"For", "!ping",
			"expected", w.ExpectedMessage,
			"got", w.ActualMessage,
		)
	}
}
