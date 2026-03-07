package main

import (
	"encoding/json"
	"net/http"
)

type MatchPayload struct {
	Result string
	Crowns int32
}

func (app *application) matchSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	var matchPayload MatchPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&matchPayload)

	if err != nil {
		http.Error(w, "Unable to parse payload", http.StatusBadRequest)
		return
	}

	// write leaderboard service

}

func (app *application) leaderboardHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromContext(r)

	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Json marshaling went wrong", http.StatusInternalServerError)
		return
	}

	w.Write(userJson)

}
