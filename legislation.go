package main

import (
	"encoding/json"
	"net/http"

	"github.com/AgoCodeBro/secret-hitler-webapp/internal/game"
)

func (gs *GameStates) discardPolicyHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name        string `json:"name"`
		PolicyIndex int    `json:"policy_index"`
	}
	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode json body", err)
		return
	}

	roomCode := r.Context().Value(gameIDKey).(string)
	g := gs.games[roomCode]

	if params.Name != g.President {
		respondWithError(w, http.StatusForbidden, "only the president can discard a policy", nil)
		return
	} else if g.CurrentPhase != game.PresidentLegislationPhase {
		respondWithError(w, http.StatusForbidden, "cannot discard policy in current phase", nil)
		return
	}

	err = g.DiscardPolicy(params.PolicyIndex)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid policy index", err)
		return
	}

	type resultJson struct {
		RoomCode          string
		RemainingPolicies []game.Policy
	}

	result := resultJson{
		RoomCode:          roomCode,
		RemainingPolicies: g.DrawnPolicies,
	}

	respondWithJson(w, http.StatusOK, result)
}

func (gs *GameStates) enactPolicyHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name        string `json:"name"`
		PolicyIndex int    `json:"policy_index"`
	}
	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode json body", err)
		return
	}

	roomCode := r.Context().Value(gameIDKey).(string)
	g := gs.games[roomCode]

	if params.Name != g.Chancelor {
		respondWithError(w, http.StatusForbidden, "only the chancelor can enact a policy", nil)
		return
	} else if g.CurrentPhase != game.ChancelorLegislationPhase {
		respondWithError(w, http.StatusForbidden, "cannot enact policy in current phase", nil)
		return
	}

	err = g.EnactPolicy(params.PolicyIndex)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid policy index", err)
		return
	}

	type resultJson struct {
		RoomCode string
	}

	result := resultJson{RoomCode: roomCode}

	respondWithJson(w, http.StatusOK, result)
}
