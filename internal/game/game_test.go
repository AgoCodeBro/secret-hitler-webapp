package game

import "testing"

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
		{"Too many players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace"}, true},
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
