package main

import "net/http"

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (app *application) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
