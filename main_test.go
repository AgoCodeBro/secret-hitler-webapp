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

func TestVoteEndpoint(t *testing.T) {
	type Tests struct {
		name        string
		playerNames []string
		votes       map[string]bool
		expectError bool
	}

	tests := []Tests{
		{
			"Valid votes",
			[]string{"Alice", "Bob", "Charlie", "David", "Eve"},
			map[string]bool{"Alice": true, "Bob": false, "Charlie": true, "David": true, "Eve": false},
			false,
		},
		{
			"Invalid player index",
			[]string{"Alice", "Bob", "Charlie", "David", "Eve"},
			map[string]bool{"Alice": true, "John": false},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameStates{games: make(map[string]*game.Game)}
			g := game.NewGame()
			g.Players = tt.playerNames
			g.President = tt.playerNames[0]
			g.Nominee = tt.playerNames[1]
			g.StartGame()
			gs.games[testGameID] = g

			type VoteRequestData struct {
				Name string `json:"name"`
				Vote bool   `json:"vote"`
			}

			gotError := false
			for name, vote := range tt.votes {
				requestData := VoteRequestData{Name: name, Vote: vote}
				jsonData, err := json.Marshal(requestData)
				if err != nil {
					t.Error(err)
				}

				req, err := http.NewRequest("POST", "/api/games/TEST/vote", bytes.NewBuffer(jsonData))
				if err != nil {
					t.Error(err)
				}

				rr := httptest.NewRecorder()

				mux := http.NewServeMux()
				mux.Handle("POST /api/games/{gameID}/vote", gs.gameMiddleware(http.HandlerFunc(gs.castVoteHandler)))
				mux.ServeHTTP(rr, req)

				if rr.Code != http.StatusOK {
					gotError = true
				}
			}

			if gotError != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, gotError)
			}
		})
	}
}

func TestDiscardEndpoint(t *testing.T) {
	type Tests struct {
		name           string
		playerNames    []string
		hand           []game.Policy
		discard        int
		expectError    bool
		expectedResult []game.Policy
	}

	tests := []Tests{
		{
			"Valid discard",
			[]string{"Alice", "Bob", "Charlie", "David", "Eve"},
			[]game.Policy{game.LiberalPolicy, game.FascistPolicy, game.LiberalPolicy},
			2,
			false,
			[]game.Policy{game.LiberalPolicy, game.LiberalPolicy},
		},
		{
			"Invalid discard index",
			[]string{"Alice", "Bob", "Charlie", "David", "Eve"},
			[]game.Policy{game.LiberalPolicy, game.FascistPolicy, game.LiberalPolicy},
			5,
			true,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameStates{games: make(map[string]*game.Game)}
			g := game.NewGame()
			g.Players = tt.playerNames
			g.President = tt.playerNames[0]
			g.Nominee = tt.playerNames[1]
			g.StartGame()
			g.DrawnPolicies = tt.hand
			g.CurrentPhase = game.PresidentLegislationPhase
			gs.games[testGameID] = g
			type DiscardRequestData struct {
				Name        string `json:"name"`
				PolicyIndex int    `json:"policy_index"`
			}
			requestData := DiscardRequestData{Name: "Alice", PolicyIndex: tt.discard}
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest("POST", "/api/games/TEST/discard", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.Handle("POST /api/games/{gameID}/discard", gs.gameMiddleware(http.HandlerFunc(gs.discardPolicyHandler)))
			mux.ServeHTTP(rr, req)

			gotError := false
			if rr.Code != http.StatusOK {
				gotError = true
			}

			if gotError != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, gotError)
			}

			if !gotError {
				for i, policy := range g.DrawnPolicies {
					if policy != tt.expectedResult[i] {
						t.Errorf("expected remaining policy %v at index %d, got %v", tt.expectedResult[i], i, policy)
					}
				}
			}
		})
	}
}

func TestEnactEndpoint(t *testing.T) {
	type Tests struct {
		name                 string
		playerNames          []string
		hand                 []game.Policy
		enact                int
		expectError          bool
		expectedLiberalCount int
		expectedFascistCount int
	}

	tests := []Tests{
		{
			"Valid enact",
			[]string{"Alice", "Bob", "Charlie", "David", "Eve"},
			[]game.Policy{game.LiberalPolicy, game.FascistPolicy},
			1,
			false,
			1,
			0,
		},
		{
			"Invalid enact index",
			[]string{"Alice", "Bob", "Charlie", "David", "Eve"},
			[]game.Policy{game.LiberalPolicy, game.FascistPolicy},
			5,
			true,
			0,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameStates{games: make(map[string]*game.Game)}
			g := game.NewGame()
			g.Players = tt.playerNames
			g.President = tt.playerNames[0]
			g.Chancelor = tt.playerNames[1]
			g.StartGame()
			g.DrawnPolicies = tt.hand
			g.CurrentPhase = game.ChancelorLegislationPhase
			gs.games[testGameID] = g
			type EnactRequestData struct {
				Name        string `json:"name"`
				PolicyIndex int    `json:"policy_index"`
			}
			requestData := EnactRequestData{Name: "Bob", PolicyIndex: tt.enact}
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest("POST", "/api/games/TEST/enact", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.Handle("POST /api/games/{gameID}/enact", gs.gameMiddleware(http.HandlerFunc(gs.enactPolicyHandler)))
			mux.ServeHTTP(rr, req)

			gotError := false
			if rr.Code != http.StatusOK {
				gotError = true
			}

			if gotError != tt.expectError {
				t.Errorf("expected error: %v, got: %v, result: %v", tt.expectError, gotError, rr.Body.String())
			}

			if !gotError {
				if g.LiberalPolicyCount != tt.expectedLiberalCount {
					t.Errorf("expected %d liberal policies enacted, got %d", tt.expectedLiberalCount, g.LiberalPolicyCount)
				}
				if g.FascistPolicyCount != tt.expectedFascistCount {
					t.Errorf("expected %d fascist policies enacted, got %d", tt.expectedFascistCount, g.FascistPolicyCount)
				}
			}
		})
	}
}
