package crater

import (
	"net/http"
	"regexp"
)

type route struct {
	pattern      *regexp.Regexp
	routeHandler http.Handler
}

type regexpHandler struct {
	getRoutes  []*route
	postRoutes []*route
}

func (h *regexpHandler) HandleGet(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.getRoutes = append(h.getRoutes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *regexpHandler) HandlePost(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.postRoutes = append(h.postRoutes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *regexpHandler) HandleStatic(pattern *regexp.Regexp, url string, fs http.FileSystem) {
	h.getRoutes = append(h.getRoutes, &route{pattern, http.StripPrefix(url, http.FileServer(fs))})
}

func (h *regexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var routes []*route
	switch r.Method {
	case "GET":
		routes = h.getRoutes
	case "POST":
		routes = h.postRoutes
	}
	for _, route := range routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.routeHandler.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

var craterRequestHandler = &regexpHandler{}
