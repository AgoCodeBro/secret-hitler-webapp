package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AgoCodeBro/secret-hitler-webapp/internal/game"
)

func (gs *GameStates) createGameHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println("A")
		respondWithError(w, http.StatusBadRequest, "failed to read request body", err)
		return
	}

	roomCode, err := gs.generateRoomCode()
	if err != nil {
		fmt.Println("B")
		respondWithError(w, http.StatusInternalServerError, "falied to generate room code", err)
		return
	}

	gs.games[roomCode] = game.NewGame()
	gs.games[roomCode].AddPlayer((params.Name))
	g := gs.games[roomCode]

	type resultJson struct {
		RoomCode string   `json:"room_code"`
		Players  []string `json:"players"`
	}

	result := resultJson{
		RoomCode: roomCode,
		Players:  g.Players,
	}

	fmt.Println(result)
	respondWithJson(w, http.StatusCreated, result)
}

func (gs *GameStates) joinGameHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name string `json:"name"`
	}

	roomCode := r.Context().Value(gameIDKey).(string)
	if roomCode == "" {
		respondWithError(w, http.StatusBadRequest, "no game code provided", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode request params when joinging game", err)
		return
	}

	if !gs.gameExists(roomCode) {
		respondWithError(w, http.StatusBadRequest, "there is no game with that room code", nil)
		return
	}

	err = gs.games[roomCode].AddPlayer(params.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to join the game", err)
		return
	}

	type resultJson struct {
		RoomCode string   `json:"room_code"`
		Players  []string `json:"players"`
	}

	result := resultJson{
		RoomCode: roomCode,
		Players:  gs.games[roomCode].Players,
	}

	respondWithJson(w, http.StatusOK, result)
}

func (gs *GameStates) startGameHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name   string `json:"name"`
		IsHost bool   `json:"is_host"`
	}

	roomCode := r.Context().Value(gameIDKey).(string)
	if roomCode == "" {
		respondWithError(w, http.StatusBadRequest, "no game code provided", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode request params when joinging game", err)
		return
	}
	if !params.IsHost {
		respondWithError(w, http.StatusForbidden, "not authorized to start game", nil)
		return
	}

	err = gs.games[roomCode].StartGame()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to start the game", err)
		return
	}

	g := gs.games[roomCode]

	type resultJson struct {
		RoomCode         string   `json:"room_code"`
		Players          []string `json:"players"`
		CurrentPresident string   `json:"current_president"`
		ElectionTracker  int      `json:"election_tracker"`
		LiberalPolicies  int      `json:"liberal_policies"`
		FascistPolicies  int      `json:"fascist_policies"`
	}

	result := resultJson{
		RoomCode:         roomCode,
		Players:          g.Players,
		CurrentPresident: g.President,
		ElectionTracker:  g.ElectionTracker,
		LiberalPolicies:  g.LiberalPolicyCount,
		FascistPolicies:  g.FascistPolicyCount,
	}

	respondWithJson(w, http.StatusOK, result)

}

func (gs *GameStates) generateRoomCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const codeLength = 4
	for {
		result := make([]byte, codeLength)
		_, err := rand.Read(result)
		if err != nil {
			return "", fmt.Errorf("failed to read random bytes: %v", err)
		}

		for i, randomByte := range result {
			letterIndex := int(randomByte) % len(charset)
			result[i] = charset[letterIndex]
		}

		if _, exist := gs.games[string(result)]; !exist {
			return string(result), nil
		}

	}
}
