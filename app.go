package crater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
)

type handlerFunc func(*Request, *Response)

// App represents you web application, which has a settings and handlers.
type App struct {
	craterRequestHandler *regexpHandler
	htmlTemplates        *craterTemplate
	settings             *Settings
	middleware           []handlerFunc
	craterRouter         *router
}

// NewApp creates a new instance of App.
// It uses default settings for application. Settings can be updated using Settings method.
func NewApp() App {
	app := App{}
	app.craterRequestHandler = newCraterHandler()
	app.htmlTemplates = &craterTemplate{}
	app.middleware = make([]handlerFunc, 0)
	app.craterRouter = new(router)
	app.settings = DefaultSettings()

	return app
}

// Settings sets settings of application.
func (app *App) Settings(settings *Settings) {
	app.settings.Update(settings)
}

// Use adds new middleware for your application.
// Middleware is a handler which is called before every request.
func (app *App) Use(handler handlerFunc) {
	app.middleware = append(app.middleware, handler)
}

// Get adds a route for HTTP GET request
func (app *App) Get(url string, handler handlerFunc) {
	requestRegexp := app.craterRouter.normalizeRoute(url)
	app.craterRequestHandler.handleGet(requestRegexp, func(w http.ResponseWriter, r *http.Request) {
		app.serveRequest(w, r, handler, requestRegexp)
	})
}

// Post adds a route for HTTP POST request.
func (app *App) Post(url string, handler handlerFunc) {
	requestRegexp := app.craterRouter.normalizeRoute(url)
	app.craterRequestHandler.handlePost(requestRegexp, func(w http.ResponseWriter, r *http.Request) {
		app.serveRequest(w, r, handler, requestRegexp)
	})
}

// Put adds a route for HTTP PUT request.
func (app *App) Put(url string, handler handlerFunc) {
	requestRegexp := app.craterRouter.normalizeRoute(url)
	app.craterRequestHandler.handlePut(requestRegexp, func(w http.ResponseWriter, r *http.Request) {
		app.serveRequest(w, r, handler, requestRegexp)
	})
}

// Delete adds a route for HTTP DELETE request.
func (app *App) Delete(url string, handler handlerFunc) {
	requestRegexp := app.craterRouter.normalizeRoute(url)
	app.craterRequestHandler.handleDelete(requestRegexp, func(w http.ResponseWriter, r *http.Request) {
		app.serveRequest(w, r, handler, requestRegexp)
	})
}

// NotFound overrides 404 status result.
// It will be called when no routes match request url.
func (app *App) NotFound(handler handlerFunc) {
	app.craterRequestHandler.notFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r, make(map[string]string))
		res := newResponse(w)
		handler(req, res)

		app.sendResponse(req, res)
	}
}

// Static handles Static files (js, css, images).
func (app *App) Static(url string) {
	app.craterRequestHandler.handleStatic(regexp.MustCompile("^"+url), url, http.Dir(app.settings.StaticPath))
}

// Listen listens on the TCP network address and then
// calls handler corresponding the request url to handle requests on incoming connections.
func (app *App) Listen(serverURL string) {
	err := app.htmlTemplates.parseFolder(app.settings.ViewsPath, app.settings.ViewExtension)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(serverURL, app.craterRequestHandler)
}

func (app *App) serveRequest(w http.ResponseWriter, r *http.Request, handler handlerFunc, requestRegexp *regexp.Regexp) {
	vars := app.craterRouter.getValues(r.URL.Path, requestRegexp)
	req := newRequest(r, vars)
	res := newResponse(w)

	if returnsResponse := app.serveMiddleware(req, res); returnsResponse {
		return
	}

	handler(req, res)

	app.sendResponse(req, res)
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
		app.redirect(res.raw, req.raw(), res.redirectUrl)
	case response_string:
		app.sendString(res.raw, res.responseString)
	}
}

func (app *App) sendJson(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", ct_JSON)
	jsonObj, err := json.Marshal(model)
	logError(err)
	fmt.Fprint(w, string(jsonObj))
}

func (app *App) sendString(w http.ResponseWriter, str string) {
	fmt.Fprint(w, str)
}

func (app *App) sendTemplate(w http.ResponseWriter, model interface{}, templateName string) {
	if templateName == "" {
		logError(fmt.Errorf("crater: TemplateName cannot be empty string"))
	}

	err := app.htmlTemplates.render(w, templateName, model)
	logError(err)
}

func (app *App) sendView(w http.ResponseWriter, model interface{}, viewName string, extension string) {
	if viewName == "" {
		logError(fmt.Errorf("crater: ViewName cannot be empty string"))
	}

	var filePath = path.Join(app.settings.ViewsPath, viewName+"."+extension)
	app.htmlTemplates.renderView(w, filePath, model)
}

func (app *App) redirect(w http.ResponseWriter, r *http.Request, url string) {
	if url == "" {
		logError(fmt.Errorf("crater: RedirectUrl cannot be empty string"))
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
