package crater

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Template struct {
	viewPath string

	tmpl *template.Template
}

func (t *Template) Parse(viewPath string, extension string) {
	t.viewPath = viewPath

	pattern := filepath.Join(viewPath, "/*."+extension)
	patternInner := filepath.Join(viewPath, "/*/*."+extension)
	tmp, _ := template.ParseGlob(pattern)
	tmp, _ = tmp.ParseGlob(patternInner)

	t.tmpl = tmp
}

func (t *Template) Render(w http.ResponseWriter, name string, data interface{}) {
	t.tmpl.ExecuteTemplate(w, name, data)
}
