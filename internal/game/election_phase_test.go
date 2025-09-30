package game

import "testing"

func TestNominateCandidate(t *testing.T) {
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

			err = g.NominateCanidate(tt.nominee)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			} else if g.NomineeIndex != tt.expectedResult {
				t.Errorf("expected nominee index: %d, got: %d", tt.expectedResult, g.NomineeIndex)
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
		expectedChencellor int
		expectedResult     bool
	}

	tests := []test{
		{"Majority yes", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, map[int]bool{0: true, 1: true, 2: false, 3: true, 4: false}, 1, true},
		{"Majority no", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, map[int]bool{0: false, 1: false, 2: true, 3: false, 4: true}, -1, false},
		{"Tie vote", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank"}, map[int]bool{0: true, 1: true, 2: false, 3: false, 4: true, 5: false}, -1, false},
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

			g.NomineeIndex = 1 // Assume Bob is nominated for simplicity

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
			if g.ChancelorIndex != tt.expectedChencellor {
				t.Errorf("expected chancellor index: %d, got: %d", tt.expectedChencellor, g.ChancelorIndex)
			}
		})
	}
}
