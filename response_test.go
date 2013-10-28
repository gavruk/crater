package crater

import "testing"

func TestRender(t *testing.T) {
	res := &Response{}
	res.Render("viewName", new(interface{}))

	if res.viewName != "viewName" {
		t.Error("viewName was not set correctly")
	}
	if res.model == nil {
		t.Error("model was not set correctly")
	}
	if res.isJson {
		t.Error("json shouldn't be set")
	}
	if res.isRedirect {
		t.Error("redirect shouldn't be set")
	}
	if res.redirectUrl != "" {
		t.Error("redirect url shouldn't be set")
	}
	if res.isHtml {
		t.Error("isHtml shouldn't be set")
	}
	if res.html != "" {
		t.Error("html shouldn't be set")
	}

}

func TestJson(t *testing.T) {
	res := &Response{}
	res.Json(new(interface{}))

	if res.model == nil {
		t.Error("model was not set correctly")
	}
	if !res.isJson {
		t.Error("json shouldn't be set")
	}
	if res.isRedirect {
		t.Error("redirect shouldn't be set")
	}
	if res.redirectUrl != "" {
		t.Error("redirect url shouldn't be set")
	}
	if res.viewName != "" {
		t.Error("viewName shouldn't be set")
	}
	if res.isHtml {
		t.Error("isHtml shouldn't be set")
	}
	if res.html != "" {
		t.Error("html shouldn't be set")
	}
}

func TestRedirect(t *testing.T) {
	res := &Response{}
	res.Redirect("redirectUrl")

	if res.redirectUrl != "redirectUrl" {
		t.Error("redirectUrl was not set correctly")
	}
	if !res.isRedirect {
		t.Error("redirect should be set")
	}
	if res.model != nil {
		t.Error("model shouldn't be set")
	}
	if res.isJson {
		t.Error("json shouldn't be set")
	}
	if res.viewName != "" {
		t.Error("viewName shouldn't be set")
	}
	if res.isHtml {
		t.Error("isHtml shouldn't be set")
	}
	if res.html != "" {
		t.Error("html shouldn't be set")
	}
}

func TestRenderString(t *testing.T) {
	res := &Response{}
	html := "<h1>text</h1>"
	res.RenderString(html)

	if !res.isHtml {
		t.Error("isHtml should be true")
	}
	if res.html != html {
		t.Error("html was not set correctly")
	}
	if res.redirectUrl != "" {
		t.Error("redirectUrl shouldn't be set")
	}
	if res.isRedirect {
		t.Error("redirect shouldn't be set")
	}
	if res.model != nil {
		t.Error("model shouldn't be set")
	}
	if res.isJson {
		t.Error("json shouldn't be set")
	}
	if res.viewName != "" {
		t.Error("viewName shouldn't be set")
	}
}

func TestResponse_CallTwoTimes_ShouldHaveLatestCallValues(t *testing.T) {
	res := &Response{}
	res.Render("viewName", new(interface{}))
	res.Redirect("redirectUrl")

	if res.redirectUrl != "redirectUrl" {
		t.Error("redirectUrl was not set correctly")
	}
	if !res.isRedirect {
		t.Error("redirect should be set")
	}
	if res.model != nil {
		t.Error("model shouldn't be set")
	}
	if res.isJson {
		t.Error("json shouldn't be set")
	}
	if res.viewName != "" {
		t.Error("viewName shouldn't be set")
	}
	if res.isHtml {
		t.Error("isHtml shouldn't be set")
	}
	if res.html != "" {
		t.Error("html shouldn't be set")
	}
}
