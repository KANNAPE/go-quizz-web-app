package handlers

import (
	"encoding/json"
	"net/http"

	"quizz-app/m/internal/chat"
	"quizz-app/m/internal/util"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, //TODO: restrict domain in prod
}

type chatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func (h *Handlers) CreateOrJoinLobby(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.sess.SetUsername(w, r, username)

	lobbyID := r.FormValue("lobbyID")
	if lobbyID != "" {
		// Join existing lobby
		if _, ok := h.store.Get(lobbyID); !ok {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, "/lobby/"+lobbyID, http.StatusSeeOther)
		return
	}

	// Create a new lobby
	l, _ := h.store.Create(username)

	http.Redirect(w, r, "/lobby/"+l.ID, http.StatusSeeOther)
}

func (h *Handlers) Lobby(w http.ResponseWriter, r *http.Request) {
	// Websocket rendering
	if websocket.IsWebSocketUpgrade(r) {
		username := h.sess.Username(r)
		if username == "" {
			http.Error(w, "Username required", http.StatusForbidden)
			return
		}
		h.handleLobbyWS(w, r)
		return
	}

	// HTML rendering
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

func (h *Handlers) handleLobbyWS(w http.ResponseWriter, r *http.Request) {
	lobbyID := mux.Vars(r)["id"]
	username := h.sess.Username(r)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &chat.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	h.chatHub.Add(lobbyID, client)

	// Reader: from client -> broadcast
	go func() {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			payload, _ := json.Marshal(chatMessage{
				Username: username,
				Text:     string(msg),
			})
			h.chatHub.Broadcast(lobbyID, payload)
		}
	}()

	// Writer: hub -> client
	go func() {
		defer conn.Close()
		for m := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, m); err != nil {
				break
			}
		}
	}()
}
