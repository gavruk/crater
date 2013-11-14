package crater

import (
	"errors"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()

	if app.craterRequestHandler == nil {
		t.Error("request handler shouldn't be nil")
	}
	if app.htmlTemplates == nil {
		t.Error("htmlTemplates shouldn't be nil")
	}
	if app.settings == nil {
		t.Error("settings shouldn't be nil")
	}
	if app.middleware == nil {
		t.Error("middleware shouldn't be nil")
	}
}

func TestGet(t *testing.T) {
	url := "/url"
	app := NewApp()
	app.Get(url, func(req *Request, res *Response) {
		res.Send("<h1>html</h1>")
	})

	if len(app.craterRequestHandler.getRoutes) != 1 {
		t.Error("get routes should have 1 hander")
	}
	err := CheckEmptyRoutesRoutes(app.craterRequestHandler.postRoutes,
		app.craterRequestHandler.putRoutes, app.craterRequestHandler.deleteRoutes)
	if err != nil {
		t.Error(err)
	}

	route := app.craterRequestHandler.getRoutes[0]
	if route.routeHandler == nil {
		t.Error("Route hander is nil")
	}
	if !route.pattern.MatchString(url) {
		t.Error("Route pattern doesn't match url")
	}
}

func TestPost(t *testing.T) {
	url := "/url"
	app := NewApp()
	app.Post(url, func(req *Request, res *Response) {
		res.Send("<h1>html</h1>")
	})

	if len(app.craterRequestHandler.postRoutes) != 1 {
		t.Error("post routes should have 1 hander")
	}
	err := CheckEmptyRoutesRoutes(app.craterRequestHandler.getRoutes,
		app.craterRequestHandler.putRoutes, app.craterRequestHandler.deleteRoutes)
	if err != nil {
		t.Error(err)
	}

	route := app.craterRequestHandler.postRoutes[0]
	if route.routeHandler == nil {
		t.Error("Route hander is nil")
	}
	if !route.pattern.MatchString(url) {
		t.Error("Route pattern doesn't match url")
	}
}

func TestPut(t *testing.T) {
	url := "/url"
	app := NewApp()
	app.Put(url, func(req *Request, res *Response) {
		res.Send("<h1>html</h1>")
	})

	if len(app.craterRequestHandler.putRoutes) != 1 {
		t.Error("put routes should have 1 hander")
	}
	err := CheckEmptyRoutesRoutes(app.craterRequestHandler.getRoutes,
		app.craterRequestHandler.postRoutes, app.craterRequestHandler.deleteRoutes)
	if err != nil {
		t.Error(err)
	}

	route := app.craterRequestHandler.putRoutes[0]
	if route.routeHandler == nil {
		t.Error("Route hander is nil")
	}
	if !route.pattern.MatchString(url) {
		t.Error("Route pattern doesn't match url")
	}
}

func TestDelete(t *testing.T) {
	url := "/url"
	app := NewApp()
	app.Delete(url, func(req *Request, res *Response) {
		res.Send("<h1>html</h1>")
	})

	if len(app.craterRequestHandler.deleteRoutes) != 1 {
		t.Error("delete routes should have 1 hander")
	}
	err := CheckEmptyRoutesRoutes(app.craterRequestHandler.getRoutes,
		app.craterRequestHandler.postRoutes, app.craterRequestHandler.putRoutes)
	if err != nil {
		t.Error(err)
	}

	route := app.craterRequestHandler.deleteRoutes[0]
	if route.routeHandler == nil {
		t.Error("Route hander is nil")
	}
	if !route.pattern.MatchString(url) {
		t.Error("Route pattern doesn't match url")
	}
}

func TestStatic(t *testing.T) {
	content := "/content"
	app := NewApp()
	app.Static("/content")

	if len(app.craterRequestHandler.getRoutes) != 1 {
		t.Error("get routes should have 1 hander")
	}
	err := CheckEmptyRoutesRoutes(app.craterRequestHandler.postRoutes,
		app.craterRequestHandler.putRoutes, app.craterRequestHandler.deleteRoutes)
	if err != nil {
		t.Error(err)
	}

	route := app.craterRequestHandler.getRoutes[0]
	if route.routeHandler == nil {
		t.Error("Route hander is nil")
	}
	if !route.pattern.MatchString(content) {
		t.Error("Route pattern doesn't match content path")
	}
}

func TestUse(t *testing.T) {
	app := NewApp()
	app.Use(func(req *Request, res *Response) {})

	if len(app.middleware) != 1 {
		t.Error("middleware should have 1 handler")
	}
}

func TestSettings(t *testing.T) {
	app := NewApp()

	settings := &Settings{
		ViewsPath:     "./folder",
		StaticPath:    "./folder",
		ViewExtension: "tmpl",
	}
	app.Settings(settings)

	if app.settings.ViewsPath != settings.ViewsPath {
		t.Error("ViewsPath was not set correctly")
	}
	if app.settings.StaticPath != settings.StaticPath {
		t.Error("StaticFilesPath was not set correctly")
	}
	if app.settings.ViewExtension != settings.ViewExtension {
		t.Error("ViewExtension was not set correctly")
	}
}

func CheckEmptyRoutesRoutes(r ...[]*route) error {
	for i := 0; i < len(r); i++ {
		if len(r[i]) != 0 {
			return errors.New("routes should be empty")
		}
	}
	return nil
}
