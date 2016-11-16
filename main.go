package main

import (
	"github.com/perthgophers/puddle/puddle"
	"os"
	"os/exec"
)

// SLACKTOKEN is the slack API token
var SLACKTOKEN string

// GITTAG Current Git Tag
var GITTAG string

func init() {
	SLACKTOKEN = os.Getenv("SLACKTOKEN")

	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		GITTAG = "NO TAG"
	} else {
		GITTAG = string(out)
	}
}

func main() {
	puddle.Run(SLACKTOKEN, GITTAG)
}
