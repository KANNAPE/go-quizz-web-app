package handlers

import (
	"log"
	"net/http"
)

func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Printf("Contact: %s %s %s",
			r.FormValue("name"), r.FormValue("email"), r.FormValue("message"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.v.Render(w, "contact.html", nil)
}
