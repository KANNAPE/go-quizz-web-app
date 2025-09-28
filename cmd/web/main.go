package main

import (
	"log"
	"net/http"

	"quizz-app/m/internal/config"
	"quizz-app/m/internal/httpx"
	"quizz-app/m/internal/lobby"
	"quizz-app/m/internal/session"
	"quizz-app/m/internal/view"
)

func main() {
	cfg := config.Load() // reads env, sets defaults

	views := view.New()             // parse+cache templates from embed
	store := lobby.NewMemoryStore() // swap later for Redis/DB
	sess := session.New(cfg.SessionKey)

	router := httpx.NewRouter(cfg, views, store, sess)

	log.Printf("listening on %s", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), router); err != nil {
		log.Fatal(err)
	}
}
