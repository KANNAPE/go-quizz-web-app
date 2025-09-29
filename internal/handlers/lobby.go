package handlers

import (
	"net/http"

	"quizz-app/m/internal/util"

	"github.com/gorilla/mux"
)

func (h *Handlers) CreateLobby(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if username == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	l, _ := h.store.Create(username)
	h.sess.SetUsername(w, r, username)
	http.Redirect(w, r, "/lobby/"+l.ID, http.StatusSeeOther)
}

func (h *Handlers) Lobby(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	l, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	username := h.sess.Username(r)
	joinURL := util.BaseURL(r) + "/lobby/" + l.ID

	data := struct {
		LobbyID  string
		Username string
		JoinURL  string
	}{
		LobbyID:  l.ID,
		Username: username,
		JoinURL:  joinURL,
	}
	h.v.Render(w, "lobby.html", data)
}
