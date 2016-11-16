package responses

import (
	"github.com/perthgophers/puddle/messagerouter"
)

// Ping responds to !ping with a pong!
func Ping(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	w.Write("Pong!")
	return nil
}

func init() {
	Handle("!ping", Ping)
}
