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
	Players        []Player
	Deck           []string
	PresidentIndex int
}

func (g *Game) resetDeck() {
	liberalCount := 6
	fascistCount := 11

	result := make([]string, (liberalCount + fascistCount))
	for i := 0; i < fascistCount; i++ {
		result[i] = "FASCIST"
	}
	for i := fascistCount; i < (liberalCount + fascistCount); i++ {
		result[i] = "LIBERAL"
	}

	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	g.Deck = result
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

func (g *Game) NominateCanidate(nominee string) (int, error) {
	for i, player := range g.Players {
		if player.Name == nominee && i != g.PresidentIndex {
			return i, nil
		}
	}
	return -1, fmt.Errorf("invalid nominee")
}
