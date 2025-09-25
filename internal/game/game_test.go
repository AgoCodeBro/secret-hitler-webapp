package game

import (
	"testing"
)

func TestAddPlayer(t *testing.T) {
	g := NewGame()
	g.AddPlayer("Alice")
	if len(g.Players) != 1 {
		t.Fatalf("expected 1 player, got %d", len(g.Players))
	}

	if g.Players[0].Name != "Alice" {
		t.Errorf("expected player name to be Alice, got %s", g.Players[0].Name)
	}
}

func TestStartGame(t *testing.T) {
	type Tests struct {
		name        string
		playerNames []string
		expectError bool
	}

	tests := []Tests{
		{"Not enough players", []string{"Alice", "Bob"}, true},
		{"Minimum players", []string{"Alice", "Bob", "Charlie", "Ago", "Kylah"}, false},
		{"Maximum players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "George", "Hector", "Ignacio", "James"}, false},
		{"Too many players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hector", "Ignacio", "James", "Kylah"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			for _, name := range tt.playerNames {
				g.AddPlayer(name)
			}
			err := g.StartGame()
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if err == nil && len(g.Players) != len(tt.playerNames) {
				t.Errorf("expected %d players, got %d", len(tt.playerNames), len(g.Players))
			}
		})
	}
}

func TestAssignRoles(t *testing.T) {
	type Tests struct {
		name             string
		playerNames      []string
		expectedLiberals int
		expectedFascists int
	}

	tests := []Tests{
		{"5 players", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, 3, 1},
		{"6 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank"}, 4, 1},
		{"7 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace"}, 4, 2},
		{"8 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank"}, 5, 2},
		{"9 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank", "Ivy"}, 5, 3},
		{"10 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank", "Ivy", "Jack"}, 6, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			for _, name := range tt.playerNames {
				g.AddPlayer(name)
			}

			err := g.StartGame()

			if err != nil {
				t.Errorf("unexpected error starting game: %v", err)
			}

			g.AssignRoles()

			var liberalCount, fascistCount, hitlerCount int
			for _, player := range g.Players {
				if player.Role == "LIBERAL" {
					liberalCount++
				} else if player.Role == "FASCIST" {
					fascistCount++
				} else if player.Role == "HITLER" {
					hitlerCount++
				} else {
					t.Errorf("unexpected role assigned: %s", player.Role)
				}
			}

			if liberalCount != tt.expectedLiberals {
				t.Errorf("expected %v liberals, got %v", tt.expectedLiberals, liberalCount)
			}
			if fascistCount != tt.expectedFascists {
				t.Errorf("expected %v fascists, got %v", tt.expectedFascists, fascistCount)
			}
			if hitlerCount != 1 {
				t.Errorf("expected 1 hitler, got %v", hitlerCount)
			}
		})
	}
}

func TestResetDeck(t *testing.T) {
	g := NewGame()
	g.resetDeck()

	if len(g.Deck) != 17 {
		t.Errorf("expected deck length of 17, got %d", len(g.Deck))
	}
	liberalCount := 0
	fascistCount := 0
	for _, card := range g.Deck {
		if card == "FASCIST" {
			fascistCount++
		} else if card == "LIBERAL" {
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

func TestNominateCanidate(t *testing.T) {
	type test struct {
		name           string
		playerNames    []string
		nominee        string
		expectedResult int
		expectError    bool
	}

	tests := []test{
		{"Valid nomination", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Bob", 1, false},
		{"Invalid nomination", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Zoe", -1, true},
		{"Self nomination", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Alice", -1, true},
		{"Valid nomination in different position", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Eve", 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			for _, name := range tt.playerNames {
				g.AddPlayer(name)
			}

			err := g.StartGame()
			if err != nil {
				t.Fatalf("unexpected error starting game: %v", err)
			}

			nomineeIndex, err := g.NominateCanidate(tt.nominee)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			} else if nomineeIndex != tt.expectedResult {
				t.Errorf("expected nominee index: %d, got: %d", tt.expectedResult, nomineeIndex)
			}

		})
	}
}
