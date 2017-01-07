package responses

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/perthgophers/puddle/messagerouter"
)

var lock = new(sync.Mutex)

// Checkout pulls & checks out a branch
func Checkout(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	var out bytes.Buffer
	words := strings.Split(cr.Text, " ")

	if len(words) < 2 {
		w.Write("I need a branch to continue with this operation")
		return nil
	}
	branch := words[1]

	cmd := exec.Command("git", "fetch", "origin")
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()

	cmd = exec.Command("git", "checkout", branch)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()

	log.Print(out.String())

	return nil
}

// Build will pull the latest git master and rebuild. It will then restart puddlebot
// Syntax: `!build <branch-name>`
// Example: `!build master`
func Build(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
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

	Checkout(cr, w)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("glide", "install")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	if stderr.Len() > 0 {
		log.Print("Error running glide up.")
		log.Print(stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	cmd = exec.Command("go", "install")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	w.Write(stdout.String())
	if stderr.Len() > 0 {
		w.WriteError("Error building:" + stderr.String())
	} else {
		w.Write(":champagne: Go install completed successfully. :champagne: ")
	}

	Restart(cr, w)
	return nil
}

func handleErr(stderr *bytes.Buffer, w messagerouter.ResponseWriter) error {
	if stderr.Len() > 0 {
		errString := stderr.String()
		stderr.Reset()
		w.WriteError(errString)
		return errors.New(errString)
	}
	return nil
}

func init() {
	Handle("!checkout", Checkout)
	Handle("!build", Build)
}
