package crater

import (
	"github.com/gavruk/crater/session"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	url := "/url"
	app := NewApp(nil)
	app.Get(url, func(req *Request, res *Response) {
		res.RenderString("<h1>html</h1>")
	})

	if len(app.craterRequestHandler.getRoutes) != 1 {
		t.Error("get routes should have 1 hander")
	}
	if len(app.craterRequestHandler.postRoutes) != 0 {
		t.Error("post routes should have no handers")
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
	app := NewApp(nil)
	app.Post(url, func(req *Request, res *Response) {
		res.RenderString("<h1>html</h1>")
	})

	if len(app.craterRequestHandler.postRoutes) != 1 {
		t.Error("post routes should have 1 hander")
	}
	if len(app.craterRequestHandler.getRoutes) != 0 {
		t.Error("get routes should have no handers")
	}
	route := app.craterRequestHandler.postRoutes[0]
	if route.routeHandler == nil {
		t.Error("Route hander is nil")
	}
	if !route.pattern.MatchString(url) {
		t.Error("Route pattern doesn't match url")
	}
}

func TestUseSessionStore(t *testing.T) {
	app := NewApp(nil)
	store := session.NewInMemorySessionStore()
	timeout := time.Minute

	sessionManager = nil
	app.UseSessionStore(store, timeout)

	if sessionManager == nil {
		t.Error("Store was not set correctly")
	}
}

func TestHandleStaticContent(t *testing.T) {
	content := "/content"
	app := NewApp(nil)
	app.HandleStaticContent("/content")

	if len(app.craterRequestHandler.postRoutes) != 0 {
		t.Error("post routes should have no handers")
	}
	if len(app.craterRequestHandler.getRoutes) != 1 {
		t.Error("get routes should have 1 hander")
	}
	route := app.craterRequestHandler.getRoutes[0]
	if route.routeHandler == nil {
		t.Error("Route hander is nil")
	}
	if !route.pattern.MatchString(content) {
		t.Error("Route pattern doesn't match content path")
	}
}
