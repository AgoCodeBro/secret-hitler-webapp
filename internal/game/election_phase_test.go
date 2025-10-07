package game

import "testing"

func TestNominateCandidate(t *testing.T) {
	type test struct {
		name           string
		playerNames    []string
		nominee        string
		expectedResult string
		expectError    bool
	}

	tests := []test{
		{"Valid nomination", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Bob", "Bob", false},
		{"Invalid nomination", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Zoe", "", true},
		{"Self nomination", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Alice", "", true},
		{"Valid nomination in different position", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Eve", "Eve", false},
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

			err = g.NominateCanidate(tt.nominee)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			} else if g.Nominee != tt.expectedResult {
				t.Errorf("expected nominee : %v, got: %v", tt.expectedResult, g.Nominee)
			}

		})
	}
}

func TestCastVote(t *testing.T) {
	type test struct {
		name        string
		playerNames []string
		votes       map[int]bool
		expectError bool
	}

	tests := []test{
		{"Valid votes", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, map[int]bool{0: true, 1: false, 2: true, 3: true, 4: false}, false},
		{"Invalid player index", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, map[int]bool{0: true, 5: false}, true},
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

			var castErrorOccurred bool
			for playerIndex, vote := range tt.votes {
				err := g.CastVote(playerIndex, vote)
				if err != nil {
					castErrorOccurred = true
				}
			}

			if castErrorOccurred != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, castErrorOccurred)
			}
		})
	}
}

func TestTallyVotes(t *testing.T) {
	type test struct {
		name               string
		playerNames        []string
		votes              map[int]bool
		expectedChencellor string
		expectedResult     bool
	}

	tests := []test{
		{"Majority yes", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, map[int]bool{0: true, 1: true, 2: false, 3: true, 4: false}, "Bob", true},
		{"Majority no", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, map[int]bool{0: false, 1: false, 2: true, 3: false, 4: true}, "", false},
		{"Tie vote", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank"}, map[int]bool{0: true, 1: true, 2: false, 3: false, 4: true, 5: false}, "", false},
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

			g.Nominee = "Bob" // Assume Bob is nominated for simplicity

			for playerIndex, vote := range tt.votes {
				err := g.CastVote(playerIndex, vote)
				if err != nil {
					t.Fatalf("unexpected error casting vote: %v", err)
				}
			}

			result := g.TallyVotes()
			if result != tt.expectedResult {
				t.Errorf("expected tally result: %v, got: %v", tt.expectedResult, result)
			}
			if g.Chancelor != tt.expectedChencellor {
				t.Errorf("expected chancellor: %v, got: %v", tt.expectedChencellor, g.Chancelor)
			}
		})
	}
}
