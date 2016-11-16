package puddle

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

var ISDEV = os.Getenv("PUDDLEDEV")

// SLACKTOKEN is the slack API token
var SLACKTOKEN string

// GITTAG is the current Git Tag
var GITTAG string

// CHANNEL is the Slack ID for channel #puddle
var CHANNEL = "C32K3QDFE"

// ISDEV Whether Puddle is in Dev or CLI mode
var ISDEV = os.Getenv("PUDDLEDEV")

var rtm *slack.RTM
var slackAPI *slack.Client

// Run Starts main Puddle process, default
// Connects to Slack & starts Slack API processing

func ErrorMessage(errText ...string) {
	for _, v := range errText {
		SendMessage(":poop: " + v + " :poop:")
	}
}
