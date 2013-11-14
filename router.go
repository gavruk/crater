package crater

import (
	"fmt"
	"regexp"
	"strings"
)

// router helps to handle complex routes
type router struct {
}

// newRouter creates new instance of router
func newRouter() *router {
	return new(router)
}

// getValues parses url and extracts route params
func (r *router) getValues(url string, pattern *regexp.Regexp) map[string]string {
	values := make(map[string]string)
	params := pattern.SubexpNames()
	if params == nil || len(params) == 0 {
		return values
	}
	routeValues := pattern.FindStringSubmatch(url)
	for i, param := range params {
		if i == 0 {
			continue
		}
		values[param] = routeValues[i]
	}
	return values
}

// normilizeRoute parses pretty route to the corresponding regexp
func (r *router) normalizeRoute(route string) *regexp.Regexp {
	pattern := regexp.MustCompile("{([^{}]*)}")
	routeValues := pattern.FindAllString(route, -1)

	if len(routeValues) == 0 {
		return r.makeRegexp(route)
	}
	for _, v := range routeValues {
		param := v[1 : len(v)-1]

		route = strings.Replace(route, v, fmt.Sprintf("(?P<%s>.*)", param), -1)
	}
	return r.makeRegexp(route)
}

func (r *router) makeRegexp(str string) *regexp.Regexp {
	return regexp.MustCompile("^" + str + "$")
}
