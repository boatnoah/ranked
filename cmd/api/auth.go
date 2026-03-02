package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string
}

// v1/register

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload UserPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&userPayload)

	if err != nil {
		http.Error(w, "Bad user request", http.StatusBadRequest)
		return
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	tok, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusBadRequest)
		return
	}

	token := Token{tok}
	tokenString, err := json.Marshal(token)
	if err != nil {
		http.Error(w, "Unable to parse token", http.StatusInternalServerError)
		return
	}

	w.Write(tokenString)
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (app *application) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
