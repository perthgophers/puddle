package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/nlopes/slack"
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
	userInfo, err := slackAPI.GetUserInfo(msg.User)
	if err != nil {
		fmt.Println(color.RedString("warning"), "No user information: ", msg.Text)
		return "", err
	}

	return userInfo.Name, nil
}

//Parse message and search admin commands and feed message
func ProcessAdminCommand(username string, msgText string, msg slack.Msg) {
	fmt.Println("ADMIN COMMAND: <"+msg.User+">", color.RedString("warning"), msg.Text)

	words := strings.Split(msgText, " ")
	cmdString := strings.TrimPrefix(words[0], "!")

	if fn, ok := adminCommands[cmdString]; ok {
		go fn(username, msgText, msg)
	} else {
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("@%s Wat?"), username))
	}
}

//Simple send message to @username: <msg>
func SendMessage(username, text string) error {
	rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("@%s: %s", username, text), CHANNEL))

	return nil
}

//Pull latest git master and rebuild, restart puddlebot
func Build(username string, msgText string, msg slack.Msg) error {
	SendMessage(username, "pulling origin/develop...")
	out, err := exec.Command("git", "pull", "origin", "develop").Output()
	SendMessage(username, string(out))
	SendMessage(username, "...Restarting...")
	cmd := exec.Command("./installandrun.sh")
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
	RegisterAdminCommand("bitcoin", BitcoinTicker)
}
