package main

import (
	"encoding/json"
	"net/http"
)

func (gs *GameStates) nominateCandidateHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		RoomCode string `json:"room_code"`
		Name     string `json:"name"`
		Nominee  string `json:"nominee"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode json body", err)
		return
	}

	game := gs.games[params.RoomCode]
	playerIndex, err := game.GetPlayerIndex(params.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to get index of player", err)
		return
	}

	if playerIndex != game.PresidentIndex {
		respondWithError(w, http.StatusForbidden, "only the president can nominate a player", nil)
		return
	}

	err = game.NominateCanidate(params.Nominee)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid nominee", err)
	}

	type resultJson struct {
		RoomCode string
		Nominee  string
	}

	result := resultJson{
		RoomCode: params.RoomCode,
		Nominee:  params.Nominee,
	}

	respondWithJson(w, http.StatusOK, result)
}

func (gs *GameStates) castVoteHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		RoomCode string `json:"room_code"`
		Name     string `json:"name"`
		Vote     bool   `json:"vote"`
	}
}
