package game

import (
	"testing"
)

func TestDrawPolicies(t *testing.T) {
	g := NewGame()
	g.resetDeck()
	expectedPolicies := []string{g.Deck[0], g.Deck[1], g.Deck[2]}
	policies := g.DrawPolicies()
	if len(policies) != 3 {
		t.Errorf("expected to draw 3 policies, got %d", len(policies))
	}
	for i, policy := range policies {
		if policy != expectedPolicies[i] {
			t.Errorf("expected policy %s, got %s", expectedPolicies[i], policy)
		}
	}
	if len(g.Deck) != 14 {
		t.Errorf("expected deck length of 14 after drawing, got %d", len(g.Deck))
	}
}

func TestDiscardPolicy(t *testing.T) {
	type test struct {
		name                    string
		policies                []string
		policyToDiscard         int
		expectedRemainingPolicy []string
		expectedError           bool
	}

	tests := []test{
		{"Valid discard", []string{"FASCIST", "LIBERAL", "FASCIST"}, 1, []string{"LIBERAL", "FASCIST"}, false},
		{"Invalid discard index", []string{"FASCIST", "LIBERAL", "FASCIST"}, 4, nil, true},
		{"Negative discard index", []string{"LIBERAL", "LIBERAL", "FASCIST"}, -1, nil, true},
		{"Discard first policy", []string{"LIBERAL", "FASCIST", "FASCIST"}, 1, []string{"FASCIST", "FASCIST"}, false},
		{"Discard last policy", []string{"FASCIST", "FASCIST", "LIBERAL"}, 3, []string{"FASCIST", "FASCIST"}, false},
		{"Empty policies", []string{}, 1, nil, true},
		{"Single policy discard", []string{"LIBERAL"}, 1, nil, true},
		{"Two policies discard first", []string{"FASCIST", "LIBERAL"}, 1, []string{"LIBERAL"}, false},
		{"Two policies discard second", []string{"FASCIST", "LIBERAL"}, 2, []string{"FASCIST"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			remainingPolicies, err := g.DiscardPolicy(tt.policies, tt.policyToDiscard)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			} else if !tt.expectedError {
				if len(remainingPolicies) != len(tt.expectedRemainingPolicy) {
					t.Errorf("expected remaining policies length: %d, got: %d", len(tt.expectedRemainingPolicy), len(remainingPolicies))
				} else {
					for i, policy := range remainingPolicies {
						if policy != tt.expectedRemainingPolicy[i] {
							t.Errorf("expected remaining policy %s at index %d, got %s", tt.expectedRemainingPolicy[i], i, policy)
						}
					}
				}
			}
		})
	}
}

func TestEnactPolicy(t *testing.T) {
	type test struct {
		name                 string
		policiesToBeEnacted  []string
		expectedLiberalCount int
		expectedFascistCount int
		expectError          bool
	}

	tests := []test{
		{"Single Liberal policy", []string{"LIBERAL"}, 1, 0, false},
		{"Single Fascist policy", []string{"FASCIST"}, 0, 1, false},
		{"Multiple Liberal policies", []string{"LIBERAL", "LIBERAL", "LIBERAL"}, 3, 0, false},
		{"Multiple Fascist policies", []string{"FASCIST", "FASCIST"}, 0, 2, false},
		{"Mixed policies", []string{"LIBERAL", "FASCIST", "LIBERAL"}, 2, 1, false},
		{"Invalid policy type", []string{"TEST"}, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			for _, policy := range tt.policiesToBeEnacted {
				err := g.EnactPolicy(policy)
				if (err != nil) != tt.expectError {
					t.Errorf("expected error: %v, got: %v", tt.expectError, err)
				}
			}
			if g.LiberalPolicyCount != tt.expectedLiberalCount {
				t.Errorf("expected %v got %v", tt.expectedLiberalCount, g.LiberalPolicyCount)
			}
			if g.FascistPolicyCount != tt.expectedFascistCount {
				t.Errorf("expected %v got %v", tt.expectedFascistCount, g.FascistPolicyCount)
			}
		})
	}
}

func TestCheckWinCondition(t *testing.T) {
	type test struct {
		name                 string
		policiesToBeEnacted  []string
		expectedWinCondition string
	}

	tests := []test{
		{"No win condition", []string{}, ""},
		{"Liberals win", []string{"LIBERAL", "LIBERAL", "LIBERAL", "LIBERAL", "LIBERAL"}, "Liberals win"},
		{"Fascists win", []string{"FASCIST", "FASCIST", "FASCIST", "FASCIST", "FASCIST", "FASCIST"}, "Fascists win"},
		{"Mixed policies no win", []string{"FASCIST", "FASCIST", "LIBERAL"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			for _, policy := range tt.policiesToBeEnacted {
				g.EnactPolicy(policy)
			}
			winCondition := g.CheckWinCondition()
			if winCondition != tt.expectedWinCondition {
				t.Errorf("expected win condition: %s, got: %s", tt.expectedWinCondition, winCondition)
			}
		})
	}
}
