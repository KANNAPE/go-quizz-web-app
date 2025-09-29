package httpx

import (
	"io/fs"
	"net/http"

	"quizz-app/m/internal/handlers"
	"quizz-app/m/web"

	"github.com/gorilla/mux"
)

func NewRouter(handler *handlers.Handlers) http.Handler {
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

	r.HandleFunc("/", handler.Home).Methods("GET")
	r.HandleFunc("/about", handler.About).Methods("GET")
	r.HandleFunc("/contact", handler.Contact).Methods("GET", "POST")
	r.HandleFunc("/lobby", handler.CreateOrJoinLobby).Methods("POST")
	r.HandleFunc("/lobby/{id}", handler.Lobby).Methods("GET")

	return r
}
