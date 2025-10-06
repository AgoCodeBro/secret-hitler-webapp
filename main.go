package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AgoCodeBro/secret-hitler-webapp/internal/game"
)

type GameStates struct {
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

	gs := GameStates{games: make(map[string]*game.Game)}

	mux.Handle("GET /api/healthz", http.HandlerFunc(readyHandler))
	mux.Handle("POST /api/games", gs.gameMiddleware(http.HandlerFunc(gs.createGameHandler)))
	mux.Handle("POST /api/games/{gameID}/join", gs.gameMiddleware(http.HandlerFunc(gs.joinGameHandler)))
	mux.Handle("POST /api/games/{gameID}/start", gs.gameMiddleware(http.HandlerFunc(gs.startGameHandler)))
	mux.Handle("POST /api/games/{gameID}/nominate", gs.gameMiddleware(http.HandlerFunc(gs.nominateCandidateHandler)))
	mux.Handle("POST /api/games/{gameID}/vote", gs.gameMiddleware(http.HandlerFunc(gs.castVoteHandler)))
	// mux.Handle("GET /api/games/{gameID}/state", gs.gameMiddleware(gs.getGameStateHandler()))

	log.Printf("Serving on port %v", port)
	log.Fatal(svr.ListenAndServe())
}

func (gs *GameStates) gameExists(code string) bool {
	_, exists := gs.games[code]
	return exists
}

func (gs *GameStates) gameMiddleware(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		roomCode := r.PathValue("gameID")
		if gs.gameExists(roomCode) {
			next.ServeHTTP(w, r)
		} else {
			respondWithError(w, http.StatusBadRequest, "game doesnt exist", nil)
		}
	}

	return http.HandlerFunc(handler)
}
