package game

import (
	"fmt"
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

type Phase string

const (
	LobbyPhase                Phase = "lobby"
	NominationPhase           Phase = "nomination"
	VotingPhase               Phase = "voting"
	PresidentLegislationPhase Phase = "president_legislation"
	ChancelorLegislationPhase Phase = "chancelor_legistlation"
	ExecutionPhase            Phase = "execution"
	GameOverPhase             Phase = "game_over"
)

type Player struct {
	Name string
	Role Role
}

type Game struct {
	Players            []Player
	CurrentPhase       Phase
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
	g.CurrentPhase = NominationPhase
}

func (g *Game) GetPlayerName(playerIndex int) (string, error) {
	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return "", fmt.Errorf("index out of range")
	}

	return g.Players[playerIndex].Name, nil
}

func (g *Game) GetPlayerIndex(playerName string) (int, error) {
	for i, player := range g.Players {
		if player.Name == playerName {
			return i, nil
		}
	}

	return 0, fmt.Errorf("player with that name not found")
}
