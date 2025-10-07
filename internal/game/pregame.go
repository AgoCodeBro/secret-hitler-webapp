package game

import (
	"fmt"
	"math/rand/v2"
)

func NewGame() *Game {
	return &Game{
		Players:      make([]string, 0, 5),
		Roles:        make(map[string]Role),
		CurrentPhase: LobbyPhase,
	}
}

func (g *Game) AddPlayer(name string) error {
	if len(g.Players) >= 6 {
		return fmt.Errorf("maximum number of players reached")
	}
	for _, player := range g.Players {
		if player == name {
			return fmt.Errorf("player with name %s already exists", name)
		}
	}
	g.Players = append(g.Players, name)
	return nil
}

func (g *Game) RemovePlayer(name string) error {
	index, err := g.GetPlayerIndex(name)
	if err != nil {
		return err
	}

	g.Players = append(g.Players[:index], g.Players[index+1:]...)
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
	g.President = g.Players[0]
	g.CurrentPhase = NominationPhase

	return nil
}

func (g *Game) AssignRoles() {
	liberalCount := (len(g.Players) / 2) + 1

	roles := make([]Role, len(g.Players))
	for i := 0; i < liberalCount; i++ {
		roles[i] = Liberal
	}
	for i := liberalCount; i < len(g.Players)-1; i++ {
		roles[i] = Fascist
	}
	roles[len(g.Players)-1] = Hitler

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	for i := range g.Players {
		fmt.Printf("Giving %v role %v\n", g.Players[i], roles[i])
		g.Roles[g.Players[i]] = roles[i]
	}
}
