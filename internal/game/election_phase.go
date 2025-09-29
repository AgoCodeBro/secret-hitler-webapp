package game

import "fmt"

func (g *Game) NominateCanidate(nominee string) (int, error) {
	for i, player := range g.Players {
		if player.Name == nominee && i != g.PresidentIndex {
			return i, nil
		}
	}
	return -1, fmt.Errorf("invalid nominee")
}

func (g *Game) CastVote(player int, vote bool) error {
	if player < 0 || player >= len(g.Players) {
		return fmt.Errorf("invalid player index")
	}
	if g.Votes == nil {
		g.Votes = make(map[int]bool)
	}
	if _, voted := g.Votes[player]; voted {
		return fmt.Errorf("player has already voted")
	}
	g.Votes[player] = vote
	return nil
}

// Returns true if the government is approved, false otherwise
func (g *Game) TallyVotes() bool {
	yesVotes := 0
	noVotes := 0
	for _, vote := range g.Votes {
		if vote {
			yesVotes++
		} else {
			noVotes++
		}
	}

	g.Votes = make(map[int]bool) // Reset votes for next round

	if yesVotes > noVotes {
		return true
	} else {
		g.PresidentIndex = (g.PresidentIndex + 1) % len(g.Players)
		return false
	}
}
