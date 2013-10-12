package crater

import "testing"

func TestRender_PassViewNameAndModel_ShouldSetThem(t *testing.T) {
	res := &Response{}
	res.Render("hello", new(interface{}))

	if res.viewName != "hello" {
		t.Error("viewName was not set correctly")
	}
	if res.model == nil {
		t.Error("model was not set correctly")
	}
}

func TestRender_JsonAndRedirectShouldBeFalse(t *testing.T) {
	res := &Response{}
	res.Render("hello", new(interface{}))

	if res.isJson {
		t.Error("isJson is true")
	}
	if res.isRedirect {
		t.Error("isRedirect is true")
	}
}

func TestJson_PassModel_ShouldSetModel(t *testing.T) {
	res := &Response{}
	res.Json(new(interface{}))

	if res.model == nil {
		t.Error("model was not set correctly")
	}
}

func TestJson_IsJsonShouldBeTrue(t *testing.T) {
	res := &Response{}
	res.Json(new(interface{}))

	if !res.isJson {
		t.Error("isJson is false")
	}
	if res.isRedirect {
		t.Error("isRedirect is true")
	}
}

func TestRedirect_ShouldSetRedirectUrl(t *testing.T) {
	res := &Response{}
	res.Redirect("hello")

	if res.redirectUrl != "hello" {
		t.Error("redirectUrl was not set")
	}
}

func TestRedirect_IsRedirectShouldBeTrue(t *testing.T) {
	res := &Response{}
	res.Redirect("hello")

	if !res.isRedirect {
		t.Error("isRedirect is false")
	}

	if res.isJson {
		t.Error("isJson is true")
	}
}
