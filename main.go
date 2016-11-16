package main

import (
	"github.com/perthgophers/puddle/messagerouter"
	"github.com/perthgophers/puddle/responses"
	"os"
	"os/exec"
)

// SLACKTOKEN is the slack API token
var SLACKTOKEN string

// GITTAG Current Git Tag
var GITTAG string

const CHANNEL = "C32K3QDFE"

func init() {
	SLACKTOKEN = os.Getenv("SLACKTOKEN")

	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		GITTAG = "NO TAG"
	} else {
		GITTAG = string(out)
	}
	GITTAG += "/"
	out, err = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err == nil {
		GITTAG += string(out)
	}
}

func main() {
	mr := messagerouter.New(SLACKTOKEN, GITTAG, CHANNEL)

	responses.Init(mr)
	mr.Run()
}
