package game

import (
	"fmt"
	"math/rand/v2"
)

func NewGame() *Game {
	return &Game{
		Players: []Player{},
	}
}

func (g *Game) AddPlayer(name string) error {
	if len(g.Players) >= 6 {
		return fmt.Errorf("maximum number of players reached")
	}
	for _, player := range g.Players {
		if player.Name == name {
			return fmt.Errorf("player with name %s already exists", name)
		}
	}
	g.Players = append(g.Players, Player{Name: name})
	return nil
}

var ErrNotEnoughPlayers = fmt.Errorf("not enough players to start the game")
var ErrTooManyPlayers = fmt.Errorf("too many players to start the game")

func (g *Game) StartGame() error {
	if len(g.Players) < 5 {
		return ErrNotEnoughPlayers
	} else if len(g.Players) > 10 {
		return ErrTooManyPlayers
	}

	g.AssignRoles()
	g.resetDeck()
	g.PresidentIndex = 0

	return nil
}

func (g *Game) AssignRoles() {
	liberalCount := (len(g.Players) / 2) + 1

	roles := make([]string, len(g.Players))
	for i := 0; i < liberalCount; i++ {
		roles[i] = "LIBERAL"
	}
	for i := liberalCount; i < len(g.Players)-1; i++ {
		roles[i] = "FASCIST"
	}
	roles[len(g.Players)-1] = "HITLER"

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	for i := range g.Players {
		g.Players[i].Role = roles[i]
	}
}
