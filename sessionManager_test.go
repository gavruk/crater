package crater

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type InMemorySessionStore struct {
	sessions map[string]*Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	store := &InMemorySessionStore{}
	store.sessions = make(map[string]*Session)

	return store
}

func (store InMemorySessionStore) Delete(id string) {
	delete(store.sessions, id)
}

func (store InMemorySessionStore) Get(id string) *Session {
	return store.sessions[id]
}

func (store InMemorySessionStore) New(id string, session *Session) {
	store.sessions[id] = session
}

func TestNewSessionManager(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour

	manager := NewSessionManager(store, timeout)

	if manager == nil {
		t.Error("SessionManager is nil")
	}
	if manager.store == nil {
		t.Error("SessionManager has no pointer to store")
	}
	if manager.timeout != timeout {
		t.Error("SessionManager has wrong timeout")
	}
}

func TestInitSession(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour
	res := newTestResponse()

	manager := NewSessionManager(store, timeout)
	session := manager.initSession(res)

	header := res.Header.Get("Set-Cookie")
	if header == "" {
		t.Error("Cookie should be set")
	}
	if session == nil {
		t.Error("initSession should return session")
	}
	if session.Value != nil {
		t.Error("New session should have a nil value")
	}
	if session.store == nil {
		t.Error("Session store shouldn't be nil")
	}
}

func TestGetSessionIdFromCookie(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour
	req := newTestRequest("GET", "localhost:8080/")
	cookieRaw := []string{"crater.SessionId=value; Expires=Fri, 11 Oct 2013 05:20:00 UTC"}
	req.Header["Cookie"] = cookieRaw

	manager := NewSessionManager(store, timeout)
	id, found := manager.getSessionIdFromCookie(req)

	if !found {
		t.Error("Existing cookie was not found")
	}
	if id != "value" {
		t.Error("Cookie value is wrong")
	}
}

func TestGetSessionIdFromCookie_WhenCookieDoesNotExist(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour
	req := newTestRequest("GET", "localhost:8080/")

	manager := NewSessionManager(store, timeout)
	id, found := manager.getSessionIdFromCookie(req)

	if found {
		t.Error("Cookie shouldn't be found")
	}
	if id != "" {
		t.Error("Cookie value should be empty string")
	}
}

func TestGetSession_NoSession_SessionShouldBeInitiated(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour
	req := newTestRequest("GET", "localhost:8080/")
	res := newTestResponse()

	manager := NewSessionManager(store, timeout)
	session := manager.GetSession(req, res)

	if session == nil {
		t.Error("Session should be initiated")
	}
}

func TestGetSession_NoSessionAndCookie_SessionShouldBeInitiated(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour
	req := newTestRequest("GET", "localhost:8080/")
	res := newTestResponse()

	manager := NewSessionManager(store, timeout)
	session := manager.GetSession(req, res)

	if session == nil {
		t.Error("Session should be initiated")
	}
	header := res.Header.Get("Set-Cookie")
	if header == "" {
		t.Error("Cookie should be set")
	}
}

func newTestRequest(method, url string) *Request {
	r, _ := http.NewRequest(method, url, nil)
	req := newRequest(r, make(map[string]string))
	return req
}

func newTestResponse() *Response {
	rec := httptest.NewRecorder()
	res := newResponse(rec)
	return res
}
