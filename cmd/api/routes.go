package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Everything is working"))
		})
		r.Post("/register", app.registerUserHandler)
		r.Post("/login", app.loginUserHandler)
		r.Route("/ranked", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/top", app.topPlayersHandler)
			r.Get("/leaderboard", app.leaderboardHandler)
			r.Post("/score", app.matchSubmissionHandler)
		})
	})

	return r
}
