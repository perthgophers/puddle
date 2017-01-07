package mud

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/perthgophers/puddle/messagerouter"
	"github.com/perthgophers/puddle/responses"
)

// Player is an struct of a player
type Player struct {
	Username     string
	CurrentPlace *Place
}

// Place is an struct of a place
type Place struct {
	Name    string
	Players map[string]*Player
	x       int
	y       int
}

var players = make(map[string]*Player)

// MudSpace is an alias for a 2D array of placess
type MudSpace [][]*Place

var mudSpace = make(MudSpace, 10)

// checkBounds returns true if the coordinates are within the array
func checkBounds(x, y int) bool {
	if x < 0 || x > len(mudSpace) || y < 0 || y > len(mudSpace[0]) {
		return false
	}
	return true
}

func (ms MudSpace) getPlace(x, y int) (*Place, bool) {
	if checkBounds(x, y) {
		return ms[x][y], true
	}
	return nil, false
}

func (ms MudSpace) moveInDirection(oldPlace *Place, direction string) (*Place, bool) {

	var newPlace *Place
	var foundPlace bool
	direction = strings.ToUpper(direction)
	switch {
	case strings.HasPrefix(direction, "N"):
		newPlace, foundPlace = ms.getPlace(oldPlace.x, oldPlace.y-1)
	case strings.HasPrefix(direction, "E"):
		newPlace, foundPlace = ms.getPlace(oldPlace.x+1, oldPlace.y)
	case strings.HasPrefix(direction, "S"):
		newPlace, foundPlace = ms.getPlace(oldPlace.x, oldPlace.y+1)
	case strings.HasPrefix(direction, "W"):
		newPlace, foundPlace = ms.getPlace(oldPlace.x-1, oldPlace.y)
	default:
		return oldPlace, false
	}

	if !foundPlace {
		return oldPlace, false
	}

	return newPlace, true
}

// InitDungeon will create a new instance of a dungeon
func InitDungeon() {
	for x := 0; x < 9; x++ {
		mudSpace[x] = make([]*Place, 10)
		for y := 0; y < 9; y++ {
			place := &Place{
				Name: fmt.Sprintf("Place %d:%d", x, y),
				x:    x,
				y:    y,
			}
			mudSpace[x][y] = place
		}
	}
}

func (p *Player) assignPlayerPosition(pl *Place) bool {
	if p.CurrentPlace != nil {
		delete(p.CurrentPlace.Players, p.Username)
	}

	if pl.Players == nil {
		pl.Players = make(map[string]*Player)
	}

	pl.Players[p.Username] = p
	p.CurrentPlace = pl

	if checkBounds(p.CurrentPlace.x, p.CurrentPlace.y) {
		return true
	}
	return false
}

// Register will add a username to the map
func Register(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	if player, foundPlayer := players[cr.Username]; foundPlayer {
		w.Write(fmt.Sprintf("You have already registered, %s.", player.Username))
		return nil
	}

	x := rand.Intn(len(mudSpace) - 1)
	y := rand.Intn(len(mudSpace[x]) - 1)
	newPlace, _ := mudSpace.getPlace(x, y)
	player := &Player{
		Username:     cr.Username,
		CurrentPlace: newPlace,
	}

	players[cr.Username] = player

	w.Write(fmt.Sprintf(
		"Welcome to Puddle Mud, %s. You are located at %d, %d.",
		cr.Username,
		player.CurrentPlace.x,
		player.CurrentPlace.y,
	))

	return nil
}

func getPlayer(username string) (*Player, error) {
	if player, ok := players[username]; ok {
		return player, nil
	}
	return nil, errors.New("User not found")
}

func (player *Player) printPosition(w messagerouter.ResponseWriter) {
	w.Write(fmt.Sprintf("You are located at %d, %d", player.CurrentPlace.x, player.CurrentPlace.y))

}

// Move will move the player in a specified direction.
// It will also prevent the player from moving off the array
func Move(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	args := strings.Split(cr.Text, " ")
	if len(args) != 2 {
		w.Write(fmt.Sprintf("%s dances some sexy moves", cr.Username))
		return nil
	}
	player, err := getPlayer(cr.Username)
	if err != nil {
		w.Write("Player not found")
		return nil
	}

	newPlace, successMove := mudSpace.moveInDirection(player.CurrentPlace, args[1])
	if !successMove {
		w.WriteError("Uncharted lands (Armadale), turn back now")
		return nil
	}
	successAssign := player.assignPlayerPosition(newPlace)
	if !successAssign {
		w.WriteError("Could not assign new position")
		return nil
	}
	player.printPosition(w)

	return nil
}

func init() {
	rand.Seed(time.Now().Unix())
	responses.Handle("%register", Register)
	responses.Handle("%move", Move)
}
