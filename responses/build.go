package responses

import (
	"fmt"
	"github.com/perthgophers/puddle/messagerouter"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var lock = new(sync.Mutex)

// Build will pull the latest git master and rebuild. It will then restart puddlebot
// Syntax: `!build <branch-name>`
// Example: `!build master`
func Build(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	w.Write("Rebuild requested. Locked.")
	lock.Lock()

	words := strings.Split(cr.Text, " ")
	branch := "master"
	if len(words) > 1 {
		branch = words[1]
	}
	w.Write(fmt.Sprintf("pulling origin/%s...", branch))

	out, err := exec.Command("git", "pull", "origin", branch).Output()
	if err != nil {
		w.WriteError("ERROR:" + err.Error())
		return err
	}
	w.Write(string(out))

	out, err = exec.Command("git", "checkout", branch).Output()
	if err != nil {
		w.Write("ERROR:" + err.Error())
		return err
	}
	w.Write(string(out))

	_, err = exec.Command("go", "install").Output()
	if err != nil {
		w.WriteError("ERROR:" + err.Error())
		return err
	}
	w.Write("...Restarting...")

	cmd := exec.Command("./run.sh")
	cmd.Start()

	os.Exit(1)
	return err
}

func init() {
	Handle("!build", Build)
}
