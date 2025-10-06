package game

import (
	"testing"
)

func TestDrawPolicies(t *testing.T) {
	g := NewGame()
	g.resetDeck()
	expectedPolicies := []Policy{g.Deck[0], g.Deck[1], g.Deck[2]}
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
		phase                   Phase
		policies                []Policy
		policyToDiscard         int
		expectedRemainingPolicy []Policy
		expectedError           bool
	}

	tests := []test{
		{"Valid discard", PresidentLegislationPhase, []Policy{FascistPolicy, LiberalPolicy, FascistPolicy}, 2, []Policy{FascistPolicy, FascistPolicy}, false},
		{"Invalid discard index", PresidentLegislationPhase, []Policy{FascistPolicy, LiberalPolicy, FascistPolicy}, 4, nil, true},
		{"Negative discard index", PresidentLegislationPhase, []Policy{FascistPolicy, LiberalPolicy, FascistPolicy}, -1, nil, true},
		{"Discard first policy", PresidentLegislationPhase, []Policy{LiberalPolicy, FascistPolicy, FascistPolicy}, 1, []Policy{FascistPolicy, FascistPolicy}, false},
		{"Discard last policy", PresidentLegislationPhase, []Policy{FascistPolicy, FascistPolicy, LiberalPolicy}, 3, []Policy{FascistPolicy, FascistPolicy}, false},
		{"Empty policies", PresidentLegislationPhase, []Policy{}, 1, nil, true},
		{"Single policy discard", PresidentLegislationPhase, []Policy{FascistPolicy}, 1, nil, true},
		{"Two policies ", PresidentLegislationPhase, []Policy{LiberalPolicy, FascistPolicy}, 1, []Policy{FascistPolicy}, true},
		{"Invalid phase", ChancelorLegislationPhase, []Policy{LiberalPolicy, FascistPolicy, FascistPolicy}, 1, []Policy{FascistPolicy}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			g.CurrentPhase = tt.phase
			g.DrawnPolicies = tt.policies
			err := g.DiscardPolicy(tt.policyToDiscard)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			} else if !tt.expectedError {
				if len(g.DrawnPolicies) != len(tt.expectedRemainingPolicy) {
					t.Errorf("expected remaining policies length: %d, got: %d", len(tt.expectedRemainingPolicy), len(g.DrawnPolicies))
				} else {
					for i, policy := range g.DrawnPolicies {
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
		drawnPolicies        []Policy
		policyToSelect       int
		expectedLiberalCount int
		expectedFascistCount int
		expectError          bool
	}

	tests := []test{
		{"Liberal policy", []Policy{LiberalPolicy, FascistPolicy}, 1, 1, 0, false},
		{"Fascist policy", []Policy{LiberalPolicy, FascistPolicy}, 2, 0, 1, false},
		{"Wrong number of policies", []Policy{LiberalPolicy, FascistPolicy, LiberalPolicy}, 1, 0, 0, true},
		{"Out of range selection", []Policy{LiberalPolicy, FascistPolicy}, 3, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			errThrown := false
			g.StartGame()
			g.AddPlayer("Ago")
			g.CurrentPhase = ChancelorLegislationPhase
			g.DrawnPolicies = tt.drawnPolicies
			err := g.EnactPolicy(tt.policyToSelect)
			if err != nil {
				errThrown = true
			}
			if errThrown != tt.expectError {
				t.Errorf("expected error: %v, got error:%v", tt.expectError, errThrown)
			}
			if g.LiberalPolicyCount != tt.expectedLiberalCount {
				t.Errorf("expected %v liberal policies got %v", tt.expectedLiberalCount, g.LiberalPolicyCount)
			}
			if g.FascistPolicyCount != tt.expectedFascistCount {
				t.Errorf("expected %v liberal policies got %v", tt.expectedFascistCount, g.FascistPolicyCount)
			}
		})
	}
}

func TestEnactTopPolicy(t *testing.T) {
	g := NewGame()
	players := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	for _, player := range players {
		g.AddPlayer(player)
	}
	g.StartGame()
	g.Deck[0] = LiberalPolicy
	g.Deck[1] = FascistPolicy

	g.EnactTopPolicy()
	if g.LiberalPolicyCount != 1 {
		t.Errorf("failed to enact correct policy from top of deck")
	}
	g.EnactTopPolicy()
	if g.FascistPolicyCount != 1 {
		t.Errorf("failed to enact correct policy from top of deck")
	}
}
