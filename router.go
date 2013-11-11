package crater

import (
	"regexp"
)

type router struct {
}

func newRouter() *router {
	return new(router)
}

func (r *router) GetValues(url string, pattern *regexp.Regexp) map[string]string {
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
