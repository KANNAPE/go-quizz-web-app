package view

import (
	"html/template"
	"net/http"
	"path"

	"quizz-app/m/web"
)

type Views struct {
	cache map[string]*template.Template
}

func New() *Views {
	v := &Views{
		cache: map[string]*template.Template{},
	}

	pages := []string{"index.html", "about.html", "contact.html", "lobby.html"}

	for _, page := range pages {
		tmpl, err := template.ParseFS(
			web.Templates,
			"templates/layout.html",
			path.Join("templates", page),
		)
		if err != nil {
			panic(err)
		}
		v.cache[page] = tmpl
	}
	return v
}

func (v *Views) Render(w http.ResponseWriter, page string, data any) {
	t, ok := v.cache[page]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}
