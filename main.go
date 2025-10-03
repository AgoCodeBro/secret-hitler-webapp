package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AgoCodeBro/secret-hitler-webapp/internal/game"
)

type GameServer struct {
	games map[string]*game.Game
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	svr := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	gs := GameServer{}

	mux.HandleFunc("GET /api/healthz", http.HandlerFunc(readyHandler))
	mux.HandleFunc("POST /api/games", http.HandlerFunc(gs.startGameHandler))
	mux.HandleFunc("POST /api/games/{gameID}/join", http.HandlerFunc(gs.joinGameHandler))

	log.Printf("Serving on port %v", port)
	log.Fatal(svr.ListenAndServe())
}

func (gs *GameServer) gameExists(code string) bool {
	_, exists := gs.games[code]
	return exists
}
