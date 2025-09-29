package game

import "fmt"

func (g *Game) DrawPolicies() []string {
	if len(g.Deck) < 3 {
		g.resetDeck()
	}

	policies := g.Deck[:3]
	g.Deck = g.Deck[3:]
	return policies
}

func (g *Game) DiscardPolicy(policies []string, policyToDiscard int) ([]string, error) {
	if len(policies) <= 1 {
		return nil, fmt.Errorf("not enough policies to discard")
	}

	if policyToDiscard < 1 || policyToDiscard > len(policies) {
		return nil, fmt.Errorf("invalid policy index to discard")
	}

	toBeRemovedIndex := policyToDiscard - 1
	remainingPolicies := make([]string, 0, len(policies)-1)
	for i, policy := range policies {
		if i != toBeRemovedIndex {
			remainingPolicies = append(remainingPolicies, policy)
		}
	}

	return remainingPolicies, nil
}

func (g *Game) EnactPolicy(policy string) error {
	switch policy {
	case "LIBERAL":
		g.LiberalPolicyCount++
	case "FASCIST":
		g.FascistPolicyCount++
	default:
		return fmt.Errorf("invalid policy type")
	}

	g.ElectionTracker = 0
	if len(g.Deck) < 3 {
		g.resetDeck()
	}
	return nil
}

func (g *Game) CheckWinCondition() string {
	if g.LiberalPolicyCount >= 5 {
		return "Liberals win"
	} else if g.FascistPolicyCount >= 6 {
		return "Fascists win"
	}

	return ""
}
