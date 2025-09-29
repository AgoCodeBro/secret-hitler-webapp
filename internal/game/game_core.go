package game

import (
	"math/rand/v2"
)

type Player struct {
	Name string
	Role string
}

type Game struct {
	Players            []Player
	Deck               []string
	PresidentIndex     int
	ChancelorIndex     int
	Votes              map[int]bool
	LiberalPolicyCount int
	FascistPolicyCount int
	ElectionTracker    int
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
