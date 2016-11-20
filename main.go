package main

import (
	"github.com/perthgophers/puddle/controllers"
	"github.com/perthgophers/puddle/messagerouter"
	"github.com/perthgophers/puddle/responses"
	"net/http"
	"os"
	"os/exec"
)

// SLACKTOKEN is the slack API token
var SLACKTOKEN string

// GITTAG Current Git Tag
var GITTAG string

const CHANNEL = "C32K3QDFE"
const SPAMCHANNEL = "C33C4MJSH"

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
	http.HandleFunc("/logs/log.html", controllers.ServeTastic)
	go http.ListenAndServe(":8080", nil)

	mr := messagerouter.New(SLACKTOKEN, GITTAG, CHANNEL, SPAMCHANNEL)

	responses.Init(mr)
	mr.Run()
}
