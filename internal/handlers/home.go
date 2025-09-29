package handlers

import (
	"net/http"

	"quizz-app/m/internal/config"
	"quizz-app/m/internal/lobby"
	"quizz-app/m/internal/session"
	"quizz-app/m/internal/view"
)

type Handlers struct {
	v     *view.Views
	store lobby.Store
	sess  *session.Manager
	cfg   config.Config
}

func New(v *view.Views, store lobby.Store, sess *session.Manager) *Handlers {
	return &Handlers{v: v, store: store, sess: sess}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	h.v.Render(w, "index.html", nil)
}
