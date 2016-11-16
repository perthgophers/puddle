package responses

import (
	"github.com/perthgophers/puddle/messagerouter"
	"os"
	"os/exec"
)

//Restart restarts Puddlebot.
func Restart(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	w.Write("...Restarting...")

	cmd := exec.Command("./run.sh")
	cmd.Start()

	os.Exit(1)
	return nil
}

func init() {
	Handle("!restart", Restart)
}
