package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/boatnoah/ranked/internal/storage"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload RegisterUserPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&userPayload)

	if err != nil {
		http.Error(w, "Bad user request", http.StatusBadRequest)
		return
	}

	var user storage.User

	user.Email = userPayload.Email
	user.Username = userPayload.Username
	err = user.Password.Set(userPayload.Password)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}

	err = app.store.UserStorage.Create(r.Context(), &user)
	if err != nil {
		http.Error(w, "unable to create user", http.StatusInternalServerError)
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
	var userPayload LoginUserPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&userPayload)

	if err != nil {
		http.Error(w, "Bad payload", http.StatusBadRequest)
		return
	}

	user, err := app.store.UserStorage.GetByEmail(r.Context(), userPayload.Email)
	if err != nil {
		http.Error(w, "No user found with those credentials", http.StatusNotFound)
		return
	}

	err = user.Password.Compare(userPayload.Password)
	if err != nil {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
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

func (app *application) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
