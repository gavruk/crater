package crater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
)

type handlerFunc func(*Request, *Response)

// App recieves settings and handles http requests
type App struct {
	craterRequestHandler *regexpHandler
	htmlTemplates        *craterTemplate
	settings             *Settings
	middleware           []handlerFunc
}

func NewApp() App {
	app := App{}
	app.craterRequestHandler = newCraterHandler()
	app.htmlTemplates = &craterTemplate{}
	app.middleware = make([]handlerFunc, 0)
	app.settings = DefaultSettings()

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

func (app *App) Use(handler handlerFunc) {
	app.middleware = append(app.middleware, handler)
}

// Get handles GET requests
func (app *App) Get(url string, handler handlerFunc) {
	app.craterRequestHandler.handleGet(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		res := newResponse(w)

		if returnsResponse := app.serveMiddleware(req, res); returnsResponse {
			return
		}

		handler(req, res)

		app.sendResponse(req, res)
	})
}

// Post handles POST requests
func (app *App) Post(url string, handler handlerFunc) {
	app.craterRequestHandler.handlePost(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		res := newResponse(w)

		if returnsResponse := app.serveMiddleware(req, res); returnsResponse {
			return
		}

		handler(req, res)

		app.sendResponse(req, res)
	})
}

// Put handles PUT requests
func (app *App) Put(url string, handler handlerFunc) {
	app.craterRequestHandler.handlePut(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		res := newResponse(w)

		if returnsResponse := app.serveMiddleware(req, res); returnsResponse {
			return
		}

		handler(req, res)

		app.sendResponse(req, res)
	})
}

// Delete handles DELETE requests
func (app *App) Delete(url string, handler handlerFunc) {
	app.craterRequestHandler.handleDelete(regexp.MustCompile("^"+url+"$"), func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		res := newResponse(w)

		if returnsResponse := app.serveMiddleware(req, res); returnsResponse {
			return
		}

		handler(req, res)

		app.sendResponse(req, res)
	})
}

// NotFound overrides 404 status result
func (app *App) NotFound(handler handlerFunc) {
	app.craterRequestHandler.notFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		res := newResponse(w)
		handler(req, res)

		app.sendResponse(req, res)
	}
}

// HandleStaticContent handles Statis Content
func (app *App) HandleStaticContent(url string) {
	app.craterRequestHandler.handleStatic(regexp.MustCompile("^"+url), url, http.Dir(app.settings.StaticFilesPath))
}

func (app *App) Listen(serverURL string) {
	err := app.htmlTemplates.parseFolder(app.settings.ViewsPath, app.settings.ViewExtension)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(serverURL, app.craterRequestHandler)
}

func (app *App) sendResponse(req *Request, res *Response) {
	switch res.responseType {
	case response_template:
		app.sendTemplate(res.raw, res.model, res.templateName)
	case response_view:
		app.sendView(res.raw, res.model, res.viewName, app.settings.ViewExtension)
	case response_json:
		app.sendJson(res.raw, res.model)
	case response_redirect:
		app.redirect(res.raw, req.raw, res.redirectUrl)
	case response_string:
		app.sendString(res.raw, res.responseString)
	}
}

func (app *App) sendJson(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", ct_JSON)
	jsonObj, _ := json.Marshal(model)
	fmt.Fprint(w, string(jsonObj))
}

func (app *App) sendString(w http.ResponseWriter, str string) {
	fmt.Fprint(w, str)
}

func (app *App) sendTemplate(w http.ResponseWriter, model interface{}, templateName string) {
	if templateName == "" {
		panic("crater: TemplateName cannot be empty string")
	}

	app.htmlTemplates.render(w, templateName, model)
}

func (app *App) sendView(w http.ResponseWriter, model interface{}, viewName string, extension string) {
	if viewName == "" {
		panic("crater: ViewName cannot be empty string")
	}

	var filePath = path.Join(app.settings.ViewsPath, viewName+"."+extension)
	app.htmlTemplates.renderView(w, filePath, model)
}

func (app *App) redirect(w http.ResponseWriter, r *http.Request, url string) {
	if url == "" {
		panic("crater: RedirectUrl cannot be empty string")
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *App) serveMiddleware(req *Request, res *Response) (returnsRespose bool) {
	returnsRespose = false
	for _, mw := range app.middleware {
		mw(req, res)
		if res.responseType != 0 {
			app.sendResponse(req, res)
			returnsRespose = true
			break
		}
	}
	return
}
