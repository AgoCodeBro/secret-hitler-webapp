package game

import (
	"fmt"
	"math/rand/v2"
)

type Player struct {
	Name string
	Role string
}

type Game struct {
	Players []Player
}

func NewGame() *Game {
	return &Game{
		Players: []Player{},
	}
}

func (g *Game) AddPlayer(name string) {
	g.Players = append(g.Players, Player{Name: name})
}

var ErrNotEnoughPlayers = fmt.Errorf("not enough players to start the game")
var ErrTooManyPlayers = fmt.Errorf("too many players to start the game")

func (g *Game) StartGame() error {
	if len(g.Players) < 5 {
		return ErrNotEnoughPlayers
	} else if len(g.Players) > 10 {
		return ErrTooManyPlayers
	}

	return nil
}

func (g *Game) AssignRoles() {
	type Role struct {
		Name      string
		SpotsLeft int
	}
	liberalCount := (len(g.Players) / 2) + 1
	fascistCount := len(g.Players) - liberalCount - 1

	roles := []Role{
		{Name: "FASCIST", SpotsLeft: fascistCount},
		{Name: "LIBERAL", SpotsLeft: liberalCount},
		{Name: "HITLER", SpotsLeft: 1},
	}

	for i := range g.Players {
		roleIndex := rand.IntN(len(roles))
		g.Players[i].Role = roles[roleIndex].Name
		roles[roleIndex].SpotsLeft--
		if roles[roleIndex].SpotsLeft == 0 {
			roles = append(roles[:roleIndex], roles[roleIndex+1:]...)
		}
	}
}
