package game

import (
	"math/rand/v2"
)

type Role string

const (
	Fascist Role = "fascist"
	Liberal Role = "liberal"
	Hitler  Role = "hitler"
)

type Policy string

const (
	FascistPolicy Policy = "facist"
	LiberalPolicy Policy = "liberal"
)

type Player struct {
	Name string
	Role Role
}

type Game struct {
	Players            []Player
	Deck               []Policy
	PresidentIndex     int
	ChancelorIndex     int
	NomineeIndex       int
	Votes              map[int]bool
	LiberalPolicyCount int
	FascistPolicyCount int
	ElectionTracker    int
}

func (g *Game) resetDeck() {
	liberalCount := 6
	fascistCount := 11

	result := make([]Policy, (liberalCount + fascistCount))
	for i := 0; i < fascistCount; i++ {
		result[i] = FascistPolicy
	}
	for i := fascistCount; i < (liberalCount + fascistCount); i++ {
		result[i] = LiberalPolicy
	}

	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	g.Deck = result
}

func (g *Game) checkWinCondition() string {
	if g.LiberalPolicyCount >= 5 {
		return "Liberals win"
	} else if g.FascistPolicyCount >= 6 {
		return "Fascists win"
	} else if g.ChancelorIndex != -1 && g.Players[g.ChancelorIndex].Role == Hitler && g.FascistPolicyCount >= 3 {
		return "Fascists win"
	}

	return ""
}

func (g *Game) StartNextRound() {
	g.ChancelorIndex = -1
	g.NomineeIndex = -1
	g.PresidentIndex = (g.PresidentIndex + 1) % len(g.Players)
}
