package crater

import (
	"encoding/json"
	"fmt"
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
	craterRequestHandler *regexpHandler
	htmlTemplates        *craterTemplate
	settings             *Settings
}

func NewApp(settings *Settings) App {
	app := App{}
	app.craterRequestHandler = newCraterHandler()
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
		req := newRequest(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := newResponse(w)
		handler(req, res)

		app.sendResponse(w, r, res)
	})
}

// Post handles POST requests
func (app App) Post(url string, handler handlerFunc) {
	app.craterRequestHandler.handlePost(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := newResponse(w)
		handler(req, res)

		app.sendResponse(w, r, res)
	})
}

// Put handles PUT requests
func (app App) Put(url string, handler handlerFunc) {
	app.craterRequestHandler.handlePut(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := newResponse(w)
		handler(req, res)

		app.sendResponse(w, r, res)
	})
}

// Delete handles DELETE requests
func (app App) Delete(url string, handler handlerFunc) {
	app.craterRequestHandler.handleDelete(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := newResponse(w)
		handler(req, res)

		app.sendResponse(w, r, res)
	})
}

// NotFound overrides 404 status result
func (app App) NotFound(handler handlerFunc) {
	app.craterRequestHandler.notFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r, sessionManager.GetSession(w, r), cookie.NewCookieManager(w, r))
		res := newResponse(w)
		handler(req, res)

		app.sendResponse(w, r, res)
	}
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

func (app App) sendResponse(w http.ResponseWriter, r *http.Request, res *Response) {
	switch res.responseType {
	case response_template:
		app.sendTemplate(w, res.model, res.templateName)
	case response_view:
		app.sendView(w, res.model, res.viewName, app.settings.ViewExtension)
	case response_json:
		app.sendJson(w, res.model)
	case response_redirect:
		app.redirect(w, r, res.redirectUrl)
	case response_string:
		app.sendString(w, res.responseString)
	}
}

func (app App) sendJson(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", ct_JSON)
	jsonObj, _ := json.Marshal(model)
	fmt.Fprint(w, string(jsonObj))
}

func (app App) sendString(w http.ResponseWriter, str string) {
	fmt.Fprint(w, str)
}

func (app App) sendTemplate(w http.ResponseWriter, model interface{}, templateName string) {
	if templateName == "" {
		panic("crater: TemplateName cannot be empty string")
	}

	app.htmlTemplates.render(w, templateName, model)
}

func (app App) sendView(w http.ResponseWriter, model interface{}, viewName string, extension string) {
	if viewName == "" {
		panic("crater: ViewName cannot be empty string")
	}

	var filePath = path.Join(app.settings.ViewsPath, viewName+"."+extension)
	app.htmlTemplates.renderView(w, filePath, model)
}

func (app App) redirect(w http.ResponseWriter, r *http.Request, url string) {
	if url == "" {
		panic("crater: RedirectUrl cannot be empty string")
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
