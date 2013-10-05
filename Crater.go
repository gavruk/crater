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
	routes []*route
}

func (h *regexpHandler) Handle(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *regexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.routeHandler.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

var craterRequestHandler = &regexpHandler{}
