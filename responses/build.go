package responses

import (
	"bytes"
	"errors"
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	out, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	currentBranch := string(out)
	w.Write("Current branch: " + currentBranch)
	w.Write("Rebuild requested. Locked.")
	lock.Lock()
	defer lock.Unlock()

	words := strings.Split(cr.Text, " ")
	branch := "master"
	if len(words) > 1 {
		branch = words[1]
	}
	w.Write(fmt.Sprintf("Selecting branch/%s...", branch))

	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	if err := handleErr(&stderr, w); err != nil {
		w.Write("Reverting back to " + currentBranch)
		cmd = exec.Command("git", "checkout", currentBranch)
		cmd.Run()

		w.Write(stdout.String())
		return err
	}
	w.Write(stdout.String())
	stdout.Reset()

	w.Write(fmt.Sprintf("Pulling origin/%s...", branch))
	cmd = exec.Command("git", "pull")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	if err := handleErr(&stderr, w); err != nil {
		w.Write("Reverting back to " + currentBranch)
		cmd = exec.Command("git", "checkout", currentBranch)
		cmd.Run()

		w.Write(stdout.String())
		return err
	}

	w.Write(stdout.String())
	stdout.Reset()

	cmd = exec.Command("go", "install")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	if err := handleErr(&stderr, w); err != nil {
		w.Write("Reverting back to " + currentBranch)
		cmd = exec.Command("git", "checkout", currentBranch)
		cmd.Run()

		w.Write(stdout.String())
		return err
	}
	w.Write("...Restarting...")
	stdout.Reset()

	cmd = exec.Command("./run.sh")
	cmd.Start()

	os.Exit(1)
	return nil
}

func handleErr(stderr *bytes.Buffer, w messagerouter.ResponseWriter) error {
	if stderr.Len() > 0 {
		errString := stderr.String()
		w.WriteError(errString)
		return errors.New(errString)
	}
	return nil
}

func init() {
	Handle("!build", Build)
}
