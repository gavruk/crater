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
	if res.responseType != response_template {
		t.Error("response type should be 'template'")
	}
}

func TestJson(t *testing.T) {
	res := &Response{}
	res.Json(new(interface{}))

	if res.model == nil {
		t.Error("model was not set correctly")
	}
	if res.responseType != response_json {
		t.Error("response type should be 'json'")
	}
}

func TestRedirect(t *testing.T) {
	res := &Response{}
	res.Redirect("redirectUrl")

	if res.redirectUrl != "redirectUrl" {
		t.Error("redirectUrl was not set correctly")
	}
	if res.responseType != response_redirect {
		t.Error("response type should be 'redirect'")
	}
}

func TestRenderString(t *testing.T) {
	res := &Response{}
	html := "<h1>text</h1>"
	res.RenderString(html)

	if res.html != html {
		t.Error("html was not set correctly")
	}
	if res.responseType != response_html {
		t.Error("response type should be 'html'")
	}
}

func TestResponse_CallTwoTimes_ShouldHaveLatestCallValues(t *testing.T) {
	res := &Response{}
	res.Render("viewName", new(interface{}))
	res.Redirect("redirectUrl")

	if res.redirectUrl != "redirectUrl" {
		t.Error("redirectUrl was not set correctly")
	}
	if res.responseType != response_redirect {
		t.Error("response type should be 'redirect'")
	}
}
