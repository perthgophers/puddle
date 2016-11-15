package puddle

import (
	"bufio"
	"fmt"
	"github.com/nlopes/slack"
	"os"
)

func RunCLI() {
	fmt.Println("Starting Puddle CLI Input...\n")
	scanner := bufio.NewScanner(os.Stdin)
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

func PrintCLI(text string) {
	fmt.Print("Puddle> ")
}
