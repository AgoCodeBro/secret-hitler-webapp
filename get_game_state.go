package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AgoCodeBro/secret-hitler-webapp/internal/game"
)

func (gs *GameStates) getGameStateHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name string
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode request body", err)
		return
	}

	g := gs.games[r.PathValue("gameID")]

	// Pregame there is no real gamestate. Give lobby info
	if g.CurrentPhase == game.LobbyPhase {
		result := struct {
			RoomCode string
			Players  []string
		}{
			RoomCode: r.PathValue("gameID"),
			Players:  g.Players,
		}

		respondWithJson(w, http.StatusOK, result)
		return
	}
	// string pointers used to allow for omit empty
	type gameState struct {
		RoomCode           string        `json:"room_code"`
		CurrentPhase       game.Phase    `json:"current_phase"`
		President          string        `json:"president"`
		Chancelor          *string       `json:"chancelor,omitempty"`
		Nominee            *string       `json:"nominee,omitempty"`
		YourHand           []game.Policy `json:"player_hand,omitempty"`
		LiberalPolicyCount int           `json:"liberal_policy_count"`
		FascistPolicyCount int           `json:"fascist_policy_count"`
		ElectionTracker    int           `json:"election_tracker"`
		Role               game.Role     `json:"role"`
		Fascists           []string      `json:"fascists,omitempty"`
		Hitler             *string       `json:"hitler,omitempty"`
	}

	role, err := g.GetPlayerRole(params.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "no role found for player", err)
		return
	}

	result := gameState{
		RoomCode:           r.PathValue("gameID"),
		CurrentPhase:       g.CurrentPhase,
		President:          g.President,
		Chancelor:          strToPtr(g.Chancelor),
		Nominee:            strToPtr(g.Nominee),
		LiberalPolicyCount: g.LiberalPolicyCount,
		FascistPolicyCount: g.FascistPolicyCount,
		ElectionTracker:    g.ElectionTracker,
		Role:               role,
	}

	if g.CurrentPhase == game.GameOverPhase {
		result.Fascists = getFascistList(g)
		hitler, err := g.GetHitler()
		if err != nil {
			log.Println(err)
		}

		result.Hitler = strToPtr(hitler)
		respondWithJson(w, http.StatusOK, result)
	}

	if g.CurrentPhase == game.PresidentLegislationPhase && params.Name == g.President {
		result.YourHand = g.DrawnPolicies
	}

	if g.CurrentPhase == game.ChancelorLegislationPhase && params.Name == g.Chancelor {
		result.YourHand = g.DrawnPolicies
	}

	if role == game.Liberal {
		respondWithJson(w, http.StatusOK, result)
		return
	}

	result.Fascists = getFascistList(g)
	if role == game.Fascist {
		hitler, err := g.GetHitler()
		if err != nil {
			log.Println(err)
		}

		result.Hitler = strToPtr(hitler)
	}

	respondWithJson(w, http.StatusOK, result)

}

func strToPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func getFascistList(g *game.Game) []string {
	var fascists []string
	for key, value := range g.Roles {
		if value == game.Fascist {
			fascists = append(fascists, key)
		}
	}

	return fascists
}
