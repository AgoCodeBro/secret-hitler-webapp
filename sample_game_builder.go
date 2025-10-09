package main

import "github.com/AgoCodeBro/secret-hitler-webapp/internal/game"

type SampleGame struct {
	g *game.Game
}

func NewSampleGame() *SampleGame {
	g := game.NewGame()
	g.Players = []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	g.CurrentPhase = game.NominationPhase
	g.President = g.Players[0]
	g.Deck = []game.Policy{game.FascistPolicy, game.LiberalPolicy, game.FascistPolicy}
	g.Roles = map[string]game.Role{"Alice": game.Liberal, "Bob": game.Fascist, "Charlie": game.Hitler, "David": game.Liberal, "Eve": game.Liberal}
	return &SampleGame{g: g}
}

func (s *SampleGame) WithPhase(phase game.Phase) {
	s.g.CurrentPhase = phase
}

func (s *SampleGame) WithPresident(president string) {
	s.g.President = president
}

func (s *SampleGame) WithChancellor(chancellor string) {
	s.g.Chancelor = chancellor
}

func (s *SampleGame) WithNominee(nominee string) {
	s.g.Nominee = nominee
}

func (s *SampleGame) WithDrawnPolicies(policies []game.Policy) {
	s.g.DrawnPolicies = policies
}

func (s *SampleGame) WithVotes(votes map[int]bool) {
	s.g.Votes = votes
}

func (s *SampleGame) WithPolicyCounts(liberal, fascist int) {
	s.g.LiberalPolicyCount = liberal
	s.g.FascistPolicyCount = fascist
}

func (s *SampleGame) WithElectionTracker(count int) {
	s.g.ElectionTracker = count
}

func (s *SampleGame) WithRoles(roles map[string]game.Role) {
	s.g.Roles = roles
}

func (s *SampleGame) WithDeck(deck []game.Policy) {
	s.g.Deck = deck
}
