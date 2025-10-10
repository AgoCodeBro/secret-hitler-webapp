package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AgoCodeBro/secret-hitler-webapp/internal/game"
)

const testGameID = "TEST"

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readyHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestCreatGameEnpoint(t *testing.T) {
	type CreateRequestData struct {
		Name string `json:"name"`
	}

	gs := GameStates{games: make(map[string]*game.Game)}

	requestData := CreateRequestData{Name: "player1"}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", "/api/games", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gs.createGameHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	type CreateResponseData struct {
		RoomCode string   `json:"room_code"`
		Players  []string `json:"players"`
	}

	var responseData CreateResponseData
	err = json.Unmarshal(rr.Body.Bytes(), &responseData)
	if err != nil {
		t.Error(err)
	}

	if responseData.RoomCode == "" {
		t.Error("expected room code to be set")
	}

	if len(responseData.Players) != 1 || responseData.Players[0] != "player1" {
		t.Errorf("expected players to contain only 'player1', got %v", responseData.Players)
	}
}

func TestJoinEndpoint(t *testing.T) {
	type JoinRequestData struct {
		Name string `json:"name"`
	}

	gs := GameStates{games: make(map[string]*game.Game)}
	newGame := game.NewGame()
	newGame.AddPlayer("player1")
	gs.games[testGameID] = newGame
	requestData := JoinRequestData{Name: "player2"}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", "/api/games/TEST/join", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gs.joinGameHandler)
	ctx := context.WithValue(req.Context(), gameIDKey, testGameID)
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v json %v",
			status, http.StatusOK, rr.Body.String())
	}

	type JoinResponseData struct {
		RoomCode string   `json:"room_code"`
		Players  []string `json:"players"`
	}

	var responseData JoinResponseData
	err = json.Unmarshal(rr.Body.Bytes(), &responseData)
	if err != nil {
		t.Error(err)
	}

	if responseData.RoomCode != testGameID {
		t.Errorf("expected room code to be 'TEST', got %v", responseData.RoomCode)
	}

	for i, player := range []string{"player1", "player2"} {
		if responseData.Players[i] != player {
			t.Errorf("expected players to contain %v at index %d, got %v", player, i, responseData.Players)
		}
	}
}

func TestStartGameEndpoint(t *testing.T) {
	type Tests struct {
		name          string
		playerNames   []string
		expectedCount int
		expectError   bool
	}

	tests := []Tests{
		{"Not enough players", []string{"Alice", "Bob"}, 2, true},
		{"Minimum players", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, 5, false},
		{"Maximum players", []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank"}, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameStates{games: make(map[string]*game.Game)}
			g := game.NewGame()
			var gotError bool
			g.Players = tt.playerNames
			gs.games[testGameID] = g

			type StartRequestData struct {
				Name   string `json:"name"`
				IsHost bool   `json:"is_host"`
			}
			requestData := StartRequestData{Name: "Alice", IsHost: true}
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest("POST", "/api/games/TEST/start", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(gs.startGameHandler)
			ctx := context.WithValue(req.Context(), gameIDKey, testGameID)
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				gotError = true
			}

			if gotError != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, gotError)
			}
			if len(g.Players) != tt.expectedCount {
				t.Errorf("expected player count: %d, got: %d", tt.expectedCount, len(g.Players))
			}
		})
	}
}

func TestNominateEndpoint(t *testing.T) {
	type Tests struct {
		name        string
		playerNames []string
		nominee     string
		expectError bool
	}

	tests := []Tests{
		{"Valid nominee", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Bob", false},
		{"Nominee is president", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Alice", true},
		{"Nominee not in game", []string{"Alice", "Bob", "Charlie", "David", "Eve"}, "Zoe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameStates{games: make(map[string]*game.Game)}
			g := game.NewGame()
			g.Players = tt.playerNames
			g.President = tt.playerNames[0]
			g.StartGame()
			gs.games[testGameID] = g

			type NominateRequestData struct {
				Name    string `json:"name"`
				Nominee string `json:"nominee"`
			}
			requestData := NominateRequestData{Name: "Alice", Nominee: tt.nominee}
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest("POST", "/api/games/TEST/nominate", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.Handle("POST /api/games/{gameID}/nominate", gs.gameMiddleware(http.HandlerFunc(gs.nominateCandidateHandler)))
			mux.ServeHTTP(rr, req)

			gotError := false
			if rr.Code != http.StatusOK {
				gotError = true
			} else {
				type NominateResponseData struct {
					RoomCode string `json:"room_code"`
					Nominee  string `json:"nominee"`
				}
				var responseData NominateResponseData
				err = json.Unmarshal(rr.Body.Bytes(), &responseData)
				if err != nil || responseData.Nominee != tt.nominee {
					gotError = true
				}
			}

			if gotError != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, gotError)
			}
		})
	}
}
