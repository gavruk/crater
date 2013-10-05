package crater

import (
	"html/template"
	"net/http"
	"regexp"
)

type handlerFunc func(*Request, *Response)

// App recieves settings and handles http requests
type App struct {
	settings Settings
}

// Settings recieves settings application
func (app *App) Settings(settings Settings) {
	app.settings = settings
}

// Get handles GET requests
func (app App) Get(url string, handler handlerFunc) {
	craterRequestHandler.Handle(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		res := &Response{}
		handler(req, res)

		t, _ := template.ParseFiles(app.settings.ViewPath + "/" + res.ViewName + ".html")
		t.Execute(w, res.Model)
	})
}
