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

type Game struct {
	Players            []string
	Roles              map[string]Role
	CurrentPhase       Phase
	Deck               []Policy
	DrawnPolicies      []Policy
	President          string
	Chancelor          string
	Nominee            string
	Votes              map[int]bool
	LiberalPolicyCount int
	FascistPolicyCount int
	ElectionTracker    int
}

func (g *Game) resetDeck() {
	const liberalCountTotal = 6
	const fascistCountTotal = 11

	liberalCount := liberalCountTotal - g.LiberalPolicyCount
	fascistCount := fascistCountTotal - g.FascistPolicyCount

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
	} else if g.Roles[g.Chancelor] == Hitler && g.FascistPolicyCount >= 3 {
		return "Fascists win"
	}

	return ""
}

func (g *Game) StartNextRound() error {
	g.Chancelor = ""
	g.Nominee = ""
	presidentIndex, err := g.GetPlayerIndex(g.President)
	if err != nil {
		return fmt.Errorf("failed to find president index: %v", err)
	}

	presidentIndex = (presidentIndex + 1) % len(g.Players)
	g.CurrentPhase = NominationPhase
	return nil
}

func (g *Game) GetPlayerName(playerIndex int) (string, error) {
	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return "", fmt.Errorf("index out of range")
	}

	return g.Players[playerIndex], nil
}

func (g *Game) GetPlayerIndex(playerName string) (int, error) {
	for i, name := range g.Players {
		if name == playerName {
			return i, nil
		}
	}

	return 0, fmt.Errorf("player with that name not found")
}
