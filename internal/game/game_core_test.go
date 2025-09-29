package game

import (
	"testing"
)

func TestResetDeck(t *testing.T) {
	g := NewGame()
	g.resetDeck()

	if len(g.Deck) != 17 {
		t.Errorf("expected deck length of 17, got %d", len(g.Deck))
	}
	liberalCount := 0
	fascistCount := 0
	for _, card := range g.Deck {
		if card == FascistPolicy {
			fascistCount++
		} else if card == LiberalPolicy {
			liberalCount++
		} else {
			t.Errorf("unexpected card in deck: %s", card)
		}
	}
	if liberalCount != 6 {
		t.Errorf("expected %v liberal cards, got %v", 6, liberalCount)
	}
	if fascistCount != 11 {
		t.Errorf("expected %v fascist cards, got %v", 11, fascistCount)
	}
}
