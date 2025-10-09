package main

import (
	"encoding/json"
	"net/http"
)

func (gs *GameStates) nominateCandidateHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name    string `json:"name"`
		Nominee string `json:"nominee"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode json body", err)
		return
	}
	roomCode := r.PathValue("gameID")
	game := gs.games[roomCode]

	if params.Name != game.President {
		respondWithError(w, http.StatusForbidden, "only the president can nominate a player", nil)
		return
	}

	err = game.NominateCanidate(params.Nominee)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid nominee", err)
		return
	}

	type resultJson struct {
		RoomCode string
		Nominee  string
	}

	result := resultJson{
		RoomCode: roomCode,
		Nominee:  params.Nominee,
	}

	respondWithJson(w, http.StatusOK, result)
}

func (gs *GameStates) castVoteHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name string `json:"name"`
		Vote bool   `json:"vote"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode json body", err)
		return
	}
	roomCode := r.PathValue("gameID")
	g := gs.games[roomCode]

	err = g.CastVote(params.Name, params.Vote)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to cast vote", err)
		return
	}

	type resultJson struct {
		RoomCode  string
		VotesCast int
		Result    bool `json:"result,omitempty"`
	}

	result := resultJson{
		RoomCode:  roomCode,
		VotesCast: len(g.Votes),
	}

	if result.VotesCast == len(g.Players) {
		result.Result = g.TallyVotes()
	}

	respondWithJson(w, http.StatusOK, result)
}
