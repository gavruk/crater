package crater

import (
	"encoding/json"
	"fmt"
	"github.com/gavruk/checker"
	"html/template"
	"net/http"
	"path"
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

		app.sendTemplate(w, res.model, res.viewName)
	})
}

// Post handles POST requests
func (app App) Post(url string, handler handlerFunc) {
	craterRequestHandler.HandlePost(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {

		req := &Request{}
		req.httpRequest = r

		res := &Response{}
		handler(req, res)

		if res.isJson {
			app.sendJson(w, res.model)
		} else {
			app.sendTemplate(w, res.model, res.viewName)
		}
	})
}

func (app App) HandleStaticFiles(url string) {
	craterRequestHandler.HandleStatic(regexp.MustCompile("^"+url), url, http.Dir(app.settings.StaticFilesPath))
}

func (app App) sendJson(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", ct_JSON)
	jsonObj, _ := json.Marshal(model)
	fmt.Fprint(w, string(jsonObj))
}

func (app App) sendTemplate(w http.ResponseWriter, model interface{}, viewName string) {
	checker.Require(viewName != "", "crater: ViewName cannot be empty string")

	t, _ := template.ParseFiles(path.Join(app.settings.ViewsPath, viewName+".html"))
	t.Execute(w, model)
}
