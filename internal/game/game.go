package game

type Player struct {
	Name string
}

type Game struct {
	Players []Player
}

func NewGame() *Game {
	return &Game{
		Players: []Player{},
	}
}

func (g *Game) AddPlayer(name string) {
	g.Players = append(g.Players, Player{Name: name})
}
