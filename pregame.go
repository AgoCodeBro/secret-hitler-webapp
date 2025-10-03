package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
)

func (gs *GameServer) startGameHandler(w http.ResponseWriter, r *http.Request) {
	var roomCode string
	var err error
	exists := true
	for exists {
		roomCode, err = generateRoomCode()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to generate room code", err)
			return
		}
		_, exists = gs.games[roomCode]
	}

	type resultJson struct {
		RoomCode string `json:"room_code"`
	}

	result := resultJson{RoomCode: roomCode}

	respondWithJson(w, http.StatusCreated, result)
}

func (gs *GameServer) joinGameHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		RoomCode string `json:"room_code"`
		Name     string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode request params when joinging game", err)
	}

	if !gs.gameExists(params.RoomCode) {
		respondWithError(w, http.StatusBadRequest, "there is no game with that room code", nil)
		return
	}

	err = gs.games[params.RoomCode].AddPlayer(params.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to join the game", err)
		return
	}

	type resultJson struct {
		RoomCode string   `json:"room_code"`
		Players  []string `json:"players"`
	}

	result := resultJson{RoomCode: params.RoomCode}
	for _, player := range gs.games[params.RoomCode].Players {
		result.Players = append(result.Players, player.Name)
	}

	respondWithJson(w, http.StatusOK, result)
}

func generateRoomCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const codeLength = 4
	result := make([]byte, codeLength)
	_, err := rand.Read(result)
	if err != nil {
		return "", fmt.Errorf("failed to read random bytes: %v", err)
	}

	for i, randomByte := range result {
		letterIndex := int(randomByte) % len(charset)
		result[i] = charset[letterIndex]
	}

	return string(result), nil
}
