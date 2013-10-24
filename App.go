package crater

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"regexp"
	"time"

	"github.com/gavruk/crater/cookie"
	"github.com/gavruk/crater/session"
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

func (app App) UseSessionStore(store session.SessionStore, timeout time.Duration) {
	sessionManager = session.NewSessionManager(store, timeout)
}

// Get handles GET requests
func (app App) Get(url string, handler handlerFunc) {
	craterRequestHandler.handleGet(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		req.init(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := &Response{}
		handler(req, res)

		if res.isJson {
			app.sendJson(w, res.model)
		} else if res.isRedirect {
			app.redirect(w, r, res.redirectUrl)
		} else if res.isHtml {
			app.sendHtml(w, res.html)
		} else {
			app.sendTemplate(w, res.model, res.viewName)
		}
	})
}

// Post handles POST requests
func (app App) Post(url string, handler handlerFunc) {
	craterRequestHandler.handlePost(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		req.init(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := &Response{}
		handler(req, res)

		if res.isJson {
			app.sendJson(w, res.model)
		} else if res.isRedirect {
			app.redirect(w, r, res.redirectUrl)
		} else if res.isHtml {
			app.sendHtml(w, res.html)
		} else {
			app.sendTemplate(w, res.model, res.viewName)
		}
	})
}

// HandleStaticContent handles Statis Content
func (app App) HandleStaticContent(url string) {
	craterRequestHandler.handleStatic(regexp.MustCompile("^"+url), url, http.Dir(app.settings.StaticFilesPath))
}

func (app App) sendJson(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", ct_JSON)
	jsonObj, _ := json.Marshal(model)
	fmt.Fprint(w, string(jsonObj))
}

func (app App) sendHtml(w http.ResponseWriter, html string) {
	fmt.Fprint(w, html)
}

func (app App) sendTemplate(w http.ResponseWriter, model interface{}, viewName string) {
	if viewName == "" {
		panic("crater: ViewName cannot be empty string")
	}

	t, _ := template.ParseFiles(path.Join(app.settings.ViewsPath, viewName+".html"))
	t.Execute(w, model)
}

func (app App) redirect(w http.ResponseWriter, r *http.Request, url string) {
	if url == "" {
		panic("crater: RedirectUrl cannot be empty string")
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
