package game

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

func (g *Game) AssignRoles() error {
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
		roleIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(roles))))
		if err != nil {
			return err
		}
		roleIndexInt := int(roleIndex.Int64())
		g.Players[i].Role = roles[roleIndexInt].Name
		roles[roleIndexInt].SpotsLeft--
		if roles[roleIndexInt].SpotsLeft == 0 {
			roles = append(roles[:roleIndexInt], roles[roleIndexInt+1:]...)
		}
	}
	return nil
}
