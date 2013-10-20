package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewCookieManager(t *testing.T) {
	r := newHttpRequest("GET", "localhost:8080/")
	w := newRecorder()

	manager := NewCookieManager(w, r)
	if manager == nil {
		t.Error("Fail to create Cookie Manager")
	}
	if manager.httpRequest == nil {
		t.Error("http request was not set")
	}
	if manager.responseWriter == nil {
		t.Error("response writer was not set")
	}
}

func TestSetCookie(t *testing.T) {
	r := newHttpRequest("GET", "localhost:8080/")
	w := newRecorder()
	expires := time.Date(2013, 10, 11, 5, 20, 0, 0, time.UTC)

	manager := NewCookieManager(w, r)
	manager.Set("name", "value", expires)

	header := manager.responseWriter.Header().Get("Set-Cookie")

	if header != "name=value; Expires=Fri, 11 Oct 2013 05:20:00 UTC" {
		t.Error("Set Cookie doesn't set correct cookie")
	}
}

func TestGetCookie(t *testing.T) {
	r := newHttpRequest("GET", "localhost:8080/")
	w := newRecorder()
	cookieRaw := []string{"name=value; Expires=Fri, 11 Oct 2013 05:20:00 UTC",
		"name2=value2; Expires=Fri, 11 Oct 2013 05:20:00 UTC"}
	r.Header["Cookie"] = cookieRaw

	manager := NewCookieManager(w, r)
	cookie, err := manager.Get("name")

	if err != nil {
		t.Error(err)
	}
	if cookie == nil {
		t.Error("Error getting cookie")
	}
	if cookie.Name != "name" {
		t.Error("Cookie name is wrong")
	}
	if cookie.Value != "value" {
		t.Error("Cookie value is wrong")
	}
}

// newHttpRequest creates a new request with a method and url
func newHttpRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func newRecorder() http.ResponseWriter {
	return httptest.NewRecorder()
}
