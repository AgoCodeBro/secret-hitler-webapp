package game

import "fmt"

type Player struct {
	Name string
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
	} else if len(g.Players) >= 6 {
		return ErrTooManyPlayers
	}
	// Additional game start logic would go here
	return nil
}
