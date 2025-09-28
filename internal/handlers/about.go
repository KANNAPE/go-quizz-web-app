package handlers

import "net/http"

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	h.v.Render(w, "about.html", nil)
}
