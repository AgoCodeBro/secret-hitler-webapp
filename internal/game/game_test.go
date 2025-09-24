package game

import "testing"

func TestAddPlayer(t *testing.T) {
	g := NewGame()
	g.AddPlayer("Alice")
	if len(g.Players) != 1 {
		t.Fatalf("expected 1 player, got %d", len(g.Players))
	}

	if g.Players[0].Name != "Alice" {
		t.Errorf("expected player name to be Alice, got %s", g.Players[0].Name)
	}
}
