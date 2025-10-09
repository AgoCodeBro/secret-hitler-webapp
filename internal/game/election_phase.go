package game

import (
	"fmt"
	"slices"
)

func (g *Game) NominateCanidate(nominee string) error {
	if slices.Contains(g.Players, nominee) && nominee != g.President {
		g.Nominee = nominee
		return nil
	}
	g.Nominee = ""
	return fmt.Errorf("invalid nominee")
}

func (g *Game) CastVote(player string, vote bool) error {
	playerIndex, err := g.GetPlayerIndex(player)
	if err != nil {
		return fmt.Errorf("player not found in game")
	}

	if playerIndex < 0 || playerIndex >= len(g.Players) {
		return fmt.Errorf("invalid player index")
	}
	if g.Votes == nil {
		g.Votes = make(map[int]bool)
	}
	if _, voted := g.Votes[playerIndex]; voted {
		return fmt.Errorf("player has already voted")
	}
	g.Votes[playerIndex] = vote
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
		g.VoteInChancellor()
		g.checkWinCondition()
		return true
	} else {
		g.ElectionTracker++
		if g.ElectionTracker >= 3 {
			g.EnactTopPolicy()
		}
		g.StartNextRound()
		return false
	}
}

func (g *Game) VoteInChancellor() {
	g.Chancelor = g.Nominee
	g.Nominee = ""
}
