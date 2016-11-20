package main

import (
	"fmt"
	"github.com/perthgophers/puddle/controllers"
	"github.com/perthgophers/puddle/messagerouter"
	"github.com/perthgophers/puddle/responses"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// SLACKTOKEN is the slack API token
var SLACKTOKEN string

// GITTAG Current Git Tag
var GITTAG string

const CHANNEL = "C32K3QDFE"
const SPAMCHANNEL = "C33C4MJSH"

// PuddleLogger is an io.Writer to add HTML to log output
type PuddleLogger struct {
	file *os.File
}

// NewPuddleLogger initialises a PuddleLogger and rotates the log file
func NewPuddleLogger() *PuddleLogger {
	pf := new(PuddleLogger)
	pf.Rotate()

	return pf
}

// Write writes len(b) bytes to the log File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
// Wraps logput in an html paragraph tag
func (pl *PuddleLogger) Write(b []byte) (n int, err error) {
	s := `<p class="log_output">` + html.EscapeString(string(b)) + `</p>`

	if pl == nil || pl.file == nil {
		return 0, os.ErrInvalid
	}
	n, e := pl.file.Write([]byte(s))

	if n < 0 {
		n = 0
	}
	if n != len(b) {
		err = io.ErrShortWrite
	}

	return n, e
}

// Close closes the log file
func (pf *PuddleLogger) Close() {
	if pf.file != nil {
		pf.file.Close()
	}
	pf.file = nil
}

// Rotate closes the opened file, and renames it to log_<datestamp>.txt for archival
// It then creates a new log.txt and opens it
func (pf *PuddleLogger) Rotate() {
	pf.Close()
	t := time.Now()
	dtSuffix := t.Format("2006_January__2_15_03_05")
	err := os.Rename("./logs/log.txt", "./logs/log_"+dtSuffix+".txt")
	if err != nil {
		fmt.Println(err)
	}
	pf.file, err = os.OpenFile("./logs/log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatalln("Unable to open log")
	}
}

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
	logger := NewPuddleLogger()

	log.SetOutput(logger)
	defer logger.Close()

	http.HandleFunc("/log", controllers.ServeTastic)
	go http.ListenAndServe(":8080", nil)

	mr := messagerouter.New(SLACKTOKEN, GITTAG, CHANNEL, SPAMCHANNEL)

	responses.Init(mr)
	mr.Run()
}
