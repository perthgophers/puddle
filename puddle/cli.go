package puddle

import (
	"bufio"
	"fmt"
	"github.com/nlopes/slack"
	"os"
)

// RunCLI Starts the command line input shell
func RunCLI() {
	fmt.Println("Starting Puddle CLI Input...\n")
	scanner := bufio.NewScanner(os.Stdin)
	SendMessage(fmt.Sprintf("... and I'm back! Git tag: %s", GITTAG))
	fmt.Print("Puddle> ")
	for scanner.Scan() {
		line := scanner.Text()
		msg := slack.Msg{
			Text: line,
		}
		ProcessMessage(msg)
		fmt.Print("Puddle> ")
	}
}
