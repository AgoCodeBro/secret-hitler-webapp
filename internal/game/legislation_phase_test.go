package game

import (
	"fmt"
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
		policies                []Policy
		policyToDiscard         int
		expectedRemainingPolicy []Policy
		expectedError           bool
	}

	tests := []test{
		{"Valid discard", []Policy{FascistPolicy, LiberalPolicy, FascistPolicy}, 2, []Policy{FascistPolicy, FascistPolicy}, false},
		{"Invalid discard index", []Policy{FascistPolicy, LiberalPolicy, FascistPolicy}, 4, nil, true},
		{"Negative discard index", []Policy{FascistPolicy, LiberalPolicy, FascistPolicy}, -1, nil, true},
		{"Discard first policy", []Policy{LiberalPolicy, FascistPolicy, FascistPolicy}, 1, []Policy{FascistPolicy, FascistPolicy}, false},
		{"Discard last policy", []Policy{FascistPolicy, FascistPolicy, LiberalPolicy}, 3, []Policy{FascistPolicy, FascistPolicy}, false},
		{"Empty policies", []Policy{}, 1, nil, true},
		{"Single policy discard", []Policy{FascistPolicy}, 1, nil, true},
		{"Two policies discard first", []Policy{LiberalPolicy, FascistPolicy}, 1, []Policy{FascistPolicy}, false},
		{"Two policies discard second", []Policy{FascistPolicy, LiberalPolicy}, 2, []Policy{FascistPolicy}, false},
		{"Two policies invalid index", []Policy{FascistPolicy, LiberalPolicy}, 3, nil, true},
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
		policiesToBeEnacted  []Policy
		expectedLiberalCount int
		expectedFascistCount int
	}

	tests := []test{
		{"Single Liberal policy", []Policy{LiberalPolicy}, 1, 0},
		{"Single Fascist policy", []Policy{FascistPolicy}, 0, 1},
		{"Multiple Liberal policies", []Policy{LiberalPolicy, LiberalPolicy}, 2, 0},
		{"Multiple Fascist policies", []Policy{FascistPolicy, FascistPolicy}, 0, 2},
		{"Mixed policies", []Policy{LiberalPolicy, FascistPolicy, LiberalPolicy}, 2, 1},
		{"No policies", []Policy{}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame()
			g.StartGame()
			g.ChancelorIndex = -1
			for _, policy := range tt.policiesToBeEnacted {
				fmt.Println(g.ChancelorIndex)
				g.EnactPolicy(policy)
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
