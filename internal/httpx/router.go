package httpx

import (
	"io/fs"
	"net/http"

	"quizz-app/m/internal/config"
	"quizz-app/m/internal/handlers"
	"quizz-app/m/internal/lobby"
	"quizz-app/m/internal/session"
	"quizz-app/m/internal/view"
	"quizz-app/m/web"

	"github.com/gorilla/mux"
)

func NewRouter(cfg config.Config, v *view.Views, store lobby.Store, sess *session.Manager) http.Handler {
	r := mux.NewRouter()

	// middleware
	r.Use(Recoverer)
	r.Use(Logger)
	r.Use(SecurityHeaders)

	// Create a filesystem rooted at "static"
	staticSub, err := fs.Sub(web.Static, "static")
	if err != nil {
		panic(err)
	}

	// Serve /static/* from that sub filesystem
	fileServer := http.FileServer(http.FS(staticSub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// handlers
	h := handlers.New(v, store, sess)

	r.HandleFunc("/", h.Home).Methods("GET")
	r.HandleFunc("/about", h.About).Methods("GET")
	r.HandleFunc("/contact", h.Contact).Methods("GET", "POST")
	r.HandleFunc("/lobby", h.CreateLobby).Methods("POST")
	r.HandleFunc("/lobby/{id}", h.Lobby).Methods("GET")

	return r
}
