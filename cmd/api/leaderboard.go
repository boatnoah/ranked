package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boatnoah/ranked/internal/leaderboard"
)

type SubmissionResponse struct {
	UserID int64
	Rank   int64
}

func (app *application) matchSubmissionHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromContext(r)
	var matchPayload leaderboard.MatchPayload

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&matchPayload)

	if err != nil {
		http.Error(w, "Unable to parse payload", http.StatusBadRequest)
		return
	}

	matchPayload.UserID = user.ID

	score, err := app.service.Submit(r.Context(), matchPayload)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	response := SubmissionResponse{UserID: user.ID, Rank: score + 1}
	jsonResponse, err := json.MarshalIndent(response, "", " ")

	if err != nil {
		http.Error(w, "Unable to parse json response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (app *application) leaderboardHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromContext(r)

	entry, err := app.service.GetPlayerRank(r.Context(), user.ID)

	if err != nil {
		http.Error(w, "Unable to retrieve player rank", http.StatusNotFound)
		return
	}

	entryJson, err := json.MarshalIndent(entry, "", " ")

	if err != nil {
		http.Error(w, "Unable marshal entry struct", http.StatusInternalServerError)
		return
	}
	w.Write(entryJson)
}
