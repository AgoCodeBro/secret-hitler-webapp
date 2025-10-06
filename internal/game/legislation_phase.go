package game

import "fmt"

func (g *Game) DrawPolicies() []Policy {
	if len(g.Deck) < 3 {
		g.resetDeck()
	}

	policies := g.Deck[:3]
	g.Deck = g.Deck[3:]
	g.DrawnPolicies = policies
	g.CurrentPhase = PresidentLegislationPhase
	return policies
}

func (g *Game) DiscardPolicy(policyToDiscard int) error {
	if g.CurrentPhase != PresidentLegislationPhase || len(g.DrawnPolicies) != 3 {
		return fmt.Errorf("invalid game state. # of policies: %v, Phase: %v", len(g.DrawnPolicies), g.CurrentPhase)
	}

	if policyToDiscard < 1 || policyToDiscard > len(g.DrawnPolicies) {
		return fmt.Errorf("invalid policy index to discard")
	}

	toBeRemovedIndex := policyToDiscard - 1
	remainingPolicies := make([]Policy, 0, len(g.DrawnPolicies)-1)
	for i, policy := range g.DrawnPolicies {
		if i != toBeRemovedIndex {
			remainingPolicies = append(remainingPolicies, policy)
		}
	}

	g.DrawnPolicies = remainingPolicies
	g.CurrentPhase = ChancelorLegislationPhase

	return nil
}

func (g *Game) EnactPolicy(policyToEnact int) error {
	if g.CurrentPhase != ChancelorLegislationPhase || len(g.DrawnPolicies) != 2 {
		return fmt.Errorf("invalid game state. # of policies: %v, phase: %v", len(g.DrawnPolicies), g.CurrentPhase)
	}

	if policyToEnact < 1 || policyToEnact > len(g.DrawnPolicies) {
		return fmt.Errorf("invalid policy index to enact")
	}

	if g.DrawnPolicies[policyToEnact-1] == LiberalPolicy {
		g.LiberalPolicyCount++
	} else {
		g.FascistPolicyCount++
	}

	g.DrawnPolicies = nil
	g.StartNextRound()
	return nil
}

func (g *Game) EnactTopPolicy() {
	policy := g.Deck[0]
	g.Deck = g.Deck[1:]
	if policy == LiberalPolicy {
		g.LiberalPolicyCount++
	} else {
		g.FascistPolicyCount++
	}
}
