package responses

import (
	"fmt"
	"github.com/perthgophers/puddle/messagerouter"
	"math/rand"
	"strings"
	"time"
)

func generateInsult() string {
	first := insult1[rand.Intn(len(insult1))]
	second := insult2[rand.Intn(len(insult2))]
	third := insult3[rand.Intn(len(insult3))]
	return fmt.Sprintf("%s %s %s", first, second, third)
}

// Insult will insult someone
func Insult(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	messageArray := strings.Split(cr.Text, " ")
	if len(messageArray) == 1 {
		insult := generateInsult()
		w.Write(fmt.Sprintf("%s, thou %s", cr.Username, insult))
	} else {
		for _, v := range messageArray[1:] {
			insult := generateInsult()
			w.Write(fmt.Sprintf("%s, thou %s", v, insult))
		}
	}

	return nil
}

func init() {
	rand.Seed(time.Now().Unix())
	Handle("!insult", Insult)
}
