package crater

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gavruk/crater/session"
	"github.com/gavruk/schema"
)

const (
	method_GET    = "GET"
	method_POST   = "POST"
	method_PUT    = "PUT"
	method_DELETE = "DELETE"

	ct_JSON              = "application/json"
	ct_FormUrlEncoded    = "application/x-www-form-urlencoded"
	ct_MultipartFormData = "multipart/form-data"
)

type route struct {
	pattern      *regexp.Regexp
	routeHandler http.Handler
}

type regexpHandler struct {
	getRoutes  []*route
	postRoutes []*route
}

func (h *regexpHandler) handleGet(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.getRoutes = append(h.getRoutes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *regexpHandler) handlePost(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.postRoutes = append(h.postRoutes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *regexpHandler) handleStatic(pattern *regexp.Regexp, url string, fs http.FileSystem) {
	h.getRoutes = append(h.getRoutes, &route{pattern, http.StripPrefix(url, http.FileServer(fs))})
}

func (h *regexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var routes []*route
	switch r.Method {
	case method_GET:
		routes = h.getRoutes
	case method_POST:
		routes = h.postRoutes
	}

	urlPath := r.URL.Path
	if urlPath == "" {
		http.NotFound(w, r)
	}
	if urlPath[0] != '/' {
		urlPath = "/" + urlPath
	}
	for _, route := range routes {
		if route.pattern.MatchString(urlPath) {
			route.routeHandler.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

var schemaDecoder = schema.NewDecoder()

var sessionManager = session.NewSessionManager(session.NewInMemorySessionStore(), time.Hour)
