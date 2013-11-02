package crater

import (
	"fmt"
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

	var tmpl, tmplInner *template.Template
	var err, errorInner error
	tmpl, err = template.ParseGlob(pattern)
	tmplInner, errorInner = tmpl.ParseGlob(patternInner)

	if err != nil && errorInner != nil {
		return fmt.Errorf("template: pattern matches no files: `%s/*.%s`", viewPath, extension)
	}

	if tmplInner != nil {
		t.ctemplate = tmplInner
	} else {
		t.ctemplate = tmpl
	}
	return nil
}

func (t *craterTemplate) render(w http.ResponseWriter, name string, data interface{}) error {
	return t.ctemplate.ExecuteTemplate(w, name, data)
}

func (t *craterTemplate) renderView(w http.ResponseWriter, path string, data interface{}) error {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	tmpl.Execute(w, data)
	return nil
}
