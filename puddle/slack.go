package puddle

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/nlopes/slack"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var lock = new(sync.Mutex)
var adminCommands = make(map[string]func(string, string, slack.Msg) error)

//Processes messages from Slack channel #puddle
//Ideally, this should be a router
func ProcessMessage(msg slack.Msg) error {
	if strings.HasPrefix(msg.Text, "!") {
		username, _ := GetUsername(msg)
		ProcessAdminCommand(username, msg.Text, msg)
	}
	return nil
}

//Retreive username from SLACK API
func GetUsername(msg slack.Msg) (string, error) {
	if ISDEV == "true" {
		return "cliuser", nil
	}
	userInfo, err := slackAPI.GetUserInfo(msg.User)
	if err != nil {
		fmt.Println(color.RedString("warning"), "No user information: ", msg.Text)
		return "", err
	}

	return userInfo.Name, nil
}

//Parse message and search admin commands and feed message
func ProcessAdminCommand(username string, msgText string, msg slack.Msg) {
	words := strings.Split(msgText, " ")
	cmdString := strings.TrimPrefix(words[0], "!")

	if fn, ok := adminCommands[cmdString]; ok {
		go fn(username, msgText, msg)
	} else {
		SendMessage(fmt.Sprintf("@%s Wat?", username))
	}
}

//Simple send message to @username: <msg>
func SendMessage(text string) error {
	if ISDEV == "true" {
		fmt.Println(">> ", text)
		return nil
	}
	rtm.SendMessage(rtm.NewOutgoingMessage(text, CHANNEL))

	return nil
}

//Pull latest git master and rebuild, restart puddlebot
func Build(username string, msgText string, msg slack.Msg) error {
	SendMessage("pulling origin/master...")
	out, err := exec.Command("git", "pull", "origin", "master").Output()
	if err != nil {
		SendMessage("ERROR:" + err.Error())
		return err
	}
	SendMessage(string(out))

	out, err = exec.Command("go", "install").Output()
	if err != nil {
		SendMessage("ERROR:" + err.Error())
		return err
	}
	SendMessage("...Restarting...")

	cmd := exec.Command("./run.sh")
	cmd.Start()

	os.Exit(1)
	return err
}

//Register admin command
func RegisterAdminCommand(cmdString string, fn func(string, string, slack.Msg) error) error {
	log.Println(fmt.Sprintf("Registering new admin command <%s>", cmdString))

	lock.Lock()
	defer lock.Unlock()

	if _, ok := adminCommands[cmdString]; !ok {
		adminCommands[cmdString] = fn
		return nil
	}

	return errors.New("Command exists")
}

func init() {
	RegisterAdminCommand("build", Build)
}
