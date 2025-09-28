package view

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"
)

// /go:embed quizz-app/m/web/templates/*.html
var templatesFS embed.FS

var staticFS embed.FS

type Views struct {
	cache  map[string]*template.Template
	static http.FileSystem
}

func New() *Views {
	v := &Views{
		cache:  map[string]*template.Template{},
		static: http.FS(staticFS),
	}
	pages := []string{"index.html", "contact.html", "lobby.html"}

	for _, page := range pages {
		t, err := template.ParseFS(templatesFS,
			"web/templates/layout.html",
			filepath.Join("web/templates", page),
		)
		if err != nil {
			panic(err)
		}
		v.cache[page] = t
	}
	return v
}

func (v *Views) Render(w http.ResponseWriter, page string, data any) {
	t, ok := v.cache[page]
	if !ok {
		http.Error(w, "template not found", 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "render error", 500)
	}
}

func (v *Views) StaticFS() http.FileSystem { return v.static }
