package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Manager struct{ store *sessions.CookieStore }

func New(key []byte) *Manager {
	cs := sessions.NewCookieStore(key)
	cs.Options.HttpOnly = true
	cs.Options.SameSite = http.SameSiteLaxMode
	return &Manager{store: cs}
}

func (m *Manager) SetUsername(w http.ResponseWriter, r *http.Request, name string) {
	s, _ := m.store.Get(r, "session")
	s.Values["username"] = name
	_ = s.Save(r, w)
}

func (m *Manager) Username(r *http.Request) string {
	s, _ := m.store.Get(r, "session")
	if v, ok := s.Values["username"].(string); ok {
		return v
	}
	return ""
}
