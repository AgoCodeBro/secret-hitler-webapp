package game

import "fmt"

func (g *Game) DrawPolicies() []Policy {
	if len(g.Deck) < 3 {
		g.resetDeck()
	}

	policies := g.Deck[:3]
	g.Deck = g.Deck[3:]
	return policies
}

func (g *Game) DiscardPolicy(policies []Policy, policyToDiscard int) ([]Policy, error) {
	if len(policies) <= 1 {
		return nil, fmt.Errorf("not enough policies to discard")
	}

	if policyToDiscard < 1 || policyToDiscard > len(policies) {
		return nil, fmt.Errorf("invalid policy index to discard")
	}

	toBeRemovedIndex := policyToDiscard - 1
	remainingPolicies := make([]Policy, 0, len(policies)-1)
	for i, policy := range policies {
		if i != toBeRemovedIndex {
			remainingPolicies = append(remainingPolicies, policy)
		}
	}

	return remainingPolicies, nil
}

func (g *Game) EnactPolicy(policy Policy) {
	if policy == LiberalPolicy {
		g.LiberalPolicyCount++
	}
	if policy == FascistPolicy {
		g.FascistPolicyCount++
	}
	g.ElectionTracker = 0
	if len(g.Deck) < 3 {
		g.resetDeck()
	}
}

func (g *Game) CheckWinCondition() string {
	if g.LiberalPolicyCount >= 5 {
		return "Liberals win"
	} else if g.FascistPolicyCount >= 6 {
		return "Fascists win"
	}

	return ""
}
