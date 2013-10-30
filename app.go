package crater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gavruk/crater/cookie"
	"github.com/gavruk/crater/session"
)

type handlerFunc func(*Request, *Response)

// App recieves settings and handles http requests
type App struct {
	craterRequestHandler *regexpHandler
	htmlTemplates        *craterTemplate
	settings             *Settings
}

func NewApp(settings *Settings) App {
	app := App{}
	app.craterRequestHandler = &regexpHandler{}
	app.htmlTemplates = &craterTemplate{}
	if settings != nil {
		app.settings = &Settings{}
		app.Settings(settings)
	} else {
		app.settings = DefaultSettings()
	}

	return app
}

// Settings recieves settings application
func (app *App) Settings(settings *Settings) {
	if settings.ViewsPath == "" {
		app.settings.ViewsPath = "."
	} else {
		app.settings.ViewsPath = settings.ViewsPath
	}
	if settings.StaticFilesPath == "" {
		app.settings.StaticFilesPath = "."
	} else {
		app.settings.StaticFilesPath = settings.StaticFilesPath
	}
	if settings.ViewExtension == "" {
		app.settings.ViewExtension = "html"
	} else {
		app.settings.ViewExtension = settings.ViewExtension
	}
}

func (app App) UseSessionStore(store session.SessionStore, timeout time.Duration) {
	sessionManager = session.NewSessionManager(store, timeout)
}

// Get handles GET requests
func (app App) Get(url string, handler handlerFunc) {
	app.craterRequestHandler.handleGet(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		req.init(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := &Response{}
		handler(req, res)

		switch res.responseType {
		case response_template:
			app.sendTemplate(w, res.model, res.viewName)
		case response_json:
			app.sendJson(w, res.model)
		case response_redirect:
			app.redirect(w, r, res.redirectUrl)
		case response_html:
			app.sendHtml(w, res.html)
		}
	})
}

// Post handles POST requests
func (app App) Post(url string, handler handlerFunc) {
	app.craterRequestHandler.handlePost(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		req.init(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := &Response{}
		handler(req, res)

		switch res.responseType {
		case response_template:
			app.sendTemplate(w, res.model, res.viewName)
		case response_json:
			app.sendJson(w, res.model)
		case response_redirect:
			app.redirect(w, r, res.redirectUrl)
		case response_html:
			app.sendHtml(w, res.html)
		}
	})
}

// HandleStaticContent handles Statis Content
func (app App) HandleStaticContent(url string) {
	app.craterRequestHandler.handleStatic(regexp.MustCompile("^"+url), url, http.Dir(app.settings.StaticFilesPath))
}

func (app App) Listen(serverURL string) {
	err := app.htmlTemplates.parseFolder(app.settings.ViewsPath, app.settings.ViewExtension)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(serverURL, app.craterRequestHandler)
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

	app.htmlTemplates.render(w, viewName, model)
}

func (app App) redirect(w http.ResponseWriter, r *http.Request, url string) {
	if url == "" {
		panic("crater: RedirectUrl cannot be empty string")
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
