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

func (t *craterTemplate) parseFolder(viewPath string, extension string) error {
	t.viewPath = viewPath

	pattern := filepath.Join(viewPath, "/*."+extension)
	patternInner := filepath.Join(viewPath, "/*/*."+extension)
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return err
	}
	tmpl, err = tmpl.ParseGlob(patternInner)
	if err != nil {
		return err
	}

	t.ctemplate = tmpl
	return nil
}

func (t *craterTemplate) render(w http.ResponseWriter, name string, data interface{}) error {
	return t.ctemplate.ExecuteTemplate(w, name, data)
}
