package mud

import (
	"fmt"
	"github.com/perthgophers/puddle/messagerouter"
	"github.com/perthgophers/puddle/responses"
)

type Player struct {
	Username string
}

type Place struct {
	Name    string
	Players []Player
	x       int
	y       int
}

var players = make(map[string]*Player)
var MudSpace = make([][]*Place, 10)

func InitDungeon() {
	for x := 0; x < 9; x++ {
		MudSpace[x] = make([]*Place, 10)
		for y := 0; y < 9; y++ {
			place := &Place{
				Name: fmt.Sprintf("Place %d:%d", x, y),
			}
			MudSpace[x][y] = place
		}
	}
}

func Register(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	if player, foundPlayer := players[cr.Username]; foundPlayer {
		w.Write(fmt.Sprintf("You have already registered, %s.", player.Username))
		return nil
	}
	player := &Player{cr.Username}
	players[cr.Username] = player

	w.Write(fmt.Sprintf("Welcome to Puddle Mud, %s", cr.Username))

	return nil
}

func init() {
	responses.Handle("%register", Register)
}
