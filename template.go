package crater

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type craterTemplate struct {
	viewPath string

	ctemplate *template.Template
}

func (t *craterTemplate) parseFolder(viewPath string, extension string) {
	t.viewPath = viewPath

	pattern := filepath.Join(viewPath, "/*."+extension)
	patternInner := filepath.Join(viewPath, "/*/*."+extension)
	tmpl, _ := template.ParseGlob(pattern)
	tmpl, _ = tmpl.ParseGlob(patternInner)

	t.ctemplate = tmpl
}

func (t *craterTemplate) render(w http.ResponseWriter, name string, data interface{}) {
	t.ctemplate.ExecuteTemplate(w, name, data)
}
