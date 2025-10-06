package game

import "testing"

func TestAddPlayer(t *testing.T) {
	type test struct {
		name           string
		playersToAdd   []string
		expectedCount  int
		expectingError bool
	}

	tests := []test{
		{"Add single player", []string{"Alice"}, 1, false},
		{"Add multiple players", []string{"Alice", "Bob", "Charlie"}, 3, false},
		{"Add duplicate player", []string{"Alice", "Bob", "Alice"}, 2, true},
		{"Add max number of players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank"}, 6, false},
		{"Add too many players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hector", "Ignacio", "James", "Kylah"}, 6, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			var err error
			gotError := false
			for _, playerName := range tt.playersToAdd {
				err = g.AddPlayer(playerName)
				if err != nil {
					gotError = true
				}
			}
			if gotError != tt.expectingError {
				t.Errorf("expected error: %v, got: %v", tt.expectingError, gotError)
			}
			if len(g.Players) != tt.expectedCount {
				t.Errorf("expected player count: %d, got: %d", tt.expectedCount, len(g.Players))
			}
		})
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
		{"Maximum players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank"}, false},
		{"Too many players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace"}, true},
		// The following tests are commented out to build the game with 5 to 6 players only first for mvp
		// They will be uncommented later once support for larger games is added

		// {"Maximum players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "George", "Hector", "Ignacio", "James"}, false},
		// {"Too many players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hector", "Ignacio", "James", "Kylah"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			gotError := false
			var err error
			for _, name := range tt.playerNames {
				err = g.AddPlayer(name)
				if err != nil {
					gotError = true
				}
			}
			err = g.StartGame()
			if err != nil {
				gotError = true
			}
			if gotError != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, gotError)
			}
			if gotError == false && len(g.Players) != len(tt.playerNames) {
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
		// The following tests are commented out to build the game with 5 to 6 players only first for mvp
		// They will be uncommented later once support for larger games is added

		// {"7 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace"}, 4, 2},
		// {"8 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank"}, 5, 2},
		// {"9 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank", "Ivy"}, 5, 3},
		// {"10 players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank", "Ivy", "Jack"}, 6, 3},
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
				if player.Role == Liberal {
					liberalCount++
				} else if player.Role == Fascist {
					fascistCount++
				} else if player.Role == Hitler {
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

func TestRemovePlayer(t *testing.T) {
	g := NewGame()
	err := g.AddPlayer("Ago")
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	err = g.AddPlayer("Kylah")
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	err = g.RemovePlayer("Kylah")
	if err != nil {
		t.Errorf("Failed to remove player: %v", err)
	}

	err = g.RemovePlayer("Kylah")
	if err == nil {
		t.Errorf("expected error")
	}

	if len(g.Players) != 1 {
		t.Errorf("expected 1 player, got %v player(s)", len(g.Players))
	}

}
