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

	if app.settings.ViewsPath == "" {
		app.settings.ViewsPath = "."
	}
	if app.settings.StaticFilesPath == "" {
		app.settings.StaticFilesPath = "."
	}
}

// Get handles GET requests
func (app App) Get(url string, handler handlerFunc) {
	craterRequestHandler.HandleGet(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {

		req := &Request{}
		req.httpRequest = r

		res := &Response{}
		handler(req, res)
		t, _ := template.ParseFiles(app.settings.ViewsPath + "/" + res.ViewName + ".html")
		t.Execute(w, res.Model)
	})
}

// Post handles POST requests
func (app App) Post(url string, handler handlerFunc) {
	craterRequestHandler.HandlePost(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {

		req := &Request{}
		req.httpRequest = r

		res := &Response{}
		handler(req, res)

		t, _ := template.ParseFiles(app.settings.ViewsPath + "/" + res.ViewName + ".html")
		t.Execute(w, res.Model)
	})
}

func (app App) HandleStaticFiles(url string) {
	craterRequestHandler.HandleStatic(regexp.MustCompile("^"+url), url, http.Dir(app.settings.StaticFilesPath))
}
