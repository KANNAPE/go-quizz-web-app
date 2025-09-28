package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var tmplCache map[string]*template.Template

// Lobby storage
type Lobby struct {
	ID       string
	HostName string
}

var lobbies = map[string]Lobby{}

// Session storage
var store = sessions.NewCookieStore([]byte("lobby-key"))

func getBaseUrl(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func initTemplates() {
	tmplCache = make(map[string]*template.Template)

	// Lisiting all pages we need to load
	pages := []string{"index.html", "about.html", "contact.html", "lobby.html"}

	for _, page := range pages {
		files := []string{
			filepath.Join("templates", "layout.html"),
			filepath.Join("templates", page),
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatalf("Error parsing template %s: %v", page, err)
		}
		tmplCache[page] = tmpl
	}
}

func renderTemplate(w http.ResponseWriter, page string, data interface{}) {
	tmpl, ok := tmplCache[page]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
	}
}

func main() {
	initTemplates()

	router := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index.html", nil)
	}).Methods("GET")

	router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "about.html", nil)
	}).Methods("GET")

	router.HandleFunc("/contact", contactPageHandler).Methods("GET", "POST")

	router.HandleFunc("/lobby", createLobbyHandler).Methods("POST")
	router.HandleFunc("/lobby/{id}", lobbyPageHandler).Methods("GET")

	fmt.Printf("Server is up and running at http://localhost")
	http.ListenAndServe(":80", router)
}

// Contact handler (not in the final version, only exists for the purpose of learning)
func contactPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		log.Printf("Contact form submitted: %s (%s): %s\n", name, email, message)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "contact.html", nil)
}

func createLobbyHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if username == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Generate unique lobby ID
	lobbyID := uuid.New().String()
	lobbies[lobbyID] = Lobby{
		ID:       lobbyID,
		HostName: username,
	}

	// Saving username in cookies
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/lobby/"+lobbyID, http.StatusSeeOther)
}

func lobbyPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lobbyID := vars["id"]

	lobby, ok := lobbies[lobbyID]
	if !ok {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	session, _ := store.Get(r, "session")
	username, _ := session.Values["username"].(string) // cast

	data := struct {
		LobbyID  string
		Username string
		JoinURL  string
	}{
		LobbyID:  lobby.ID,
		Username: username,
		JoinURL:  getBaseUrl(r) + "/lobby/" + lobby.ID,
	}

	renderTemplate(w, "lobby.html", data)
}
