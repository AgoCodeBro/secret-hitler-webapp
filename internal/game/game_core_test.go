package game

import (
	"testing"
)

func TestResetDeck(t *testing.T) {
	g := NewGame()
	g.resetDeck()

	if len(g.Deck) != 17 {
		t.Errorf("expected deck length of 17, got %d", len(g.Deck))
	}
	liberalCount := 0
	fascistCount := 0
	for _, card := range g.Deck {
		if card == FascistPolicy {
			fascistCount++
		} else if card == LiberalPolicy {
			liberalCount++
		} else {
			t.Errorf("unexpected card in deck: %s", card)
		}
	}
	if liberalCount != 6 {
		t.Errorf("expected %v liberal cards, got %v", 6, liberalCount)
	}
	if fascistCount != 11 {
		t.Errorf("expected %v fascist cards, got %v", 11, fascistCount)
	}
}

func TestCheckWinConditionFromPolicy(t *testing.T) {
	type test struct {
		name                 string
		liberalCount         int
		fascistCount         int
		expectedWinCondition string
	}

	tests := []test{
		{"No win condition", 0, 0, ""},
		{"Liberals win", 5, 3, "Liberals win"},
		{"Fascists win", 2, 6, "Fascists win"},
		{"Mixed policies no win", 2, 3, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			players := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
			for _, player := range players {
				g.AddPlayer(player)
			}
			g.StartGame()

			g.LiberalPolicyCount = tt.liberalCount
			g.FascistPolicyCount = tt.fascistCount
			winCondition := g.checkWinCondition()
			if winCondition != tt.expectedWinCondition {
				t.Errorf("expected win condition: %s, got: %s", tt.expectedWinCondition, winCondition)
			}
		})
	}
}

func TestWinConditionFromElection(t *testing.T) {
	g := NewGame()
	g.Players = []Player{
		{Name: "Ago", Role: Fascist},
		{Name: "Kylah", Role: Liberal},
		{Name: "Jaren", Role: Hitler},
	}
	for i := 0; i < 5; i++ {
		g.FascistPolicyCount = i
		for j := range g.Players {
			if i < 3 {
				g.ChancelorIndex = j
				result := g.checkWinCondition()
				if result != "" {
					t.Errorf("winner declared when there should not be one: %v", result)
				}
			} else {
				g.ChancelorIndex = j
				result := g.checkWinCondition()
				if j == 2 {
					if result != "Fascists win" {
						t.Errorf("fascist should have won: %v", result)
					}
				} else if result != "" {
					t.Errorf("winner declared when there should not be one: %v", result)
				}
			}

		}
	}

}
