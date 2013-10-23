package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// -----------------------------
// In-Memory Session Store
// -----------------------------

func TestNewInMemorySessionStore(t *testing.T) {
	store := NewInMemorySessionStore()

	if store == nil {
		t.Error("InMemorySessionStore is nil")
	}

	if store.sessions == nil {
		t.Error("session shouldn't be nil")
	}
}

func TestNew(t *testing.T) {
	store := NewInMemorySessionStore()
	id := "abc"
	session := &Session{
		Id:      "id",
		Value:   42,
		Expires: time.Date(2013, 10, 11, 5, 20, 0, 0, time.UTC),
	}
	store.New(id, session)

	if len(store.sessions) != 1 {
		t.Error("sessions count in the store should be 1")
	}

	storeSession := store.sessions[id]
	if storeSession == nil {
		t.Error("Session was not found by id")
	}
	if storeSession.Id != session.Id {
		t.Error("Wrong session id")
	}
	if storeSession.Value.(int) != session.Value.(int) {
		t.Error("Wrong session value")
	}
	if storeSession.Expires != session.Expires {
		t.Error("Wrong session expires date")
	}
}

func TestGet(t *testing.T) {
	store := NewInMemorySessionStore()
	id := "abc"
	session := &Session{
		Id:      "id",
		Value:   42,
		Expires: time.Date(2013, 10, 11, 5, 20, 0, 0, time.UTC),
	}
	store.New(id, session)

	storeSession := store.Get(id)
	if storeSession == nil {
		t.Error("Session was not found by id")
	}
	if storeSession.Id != session.Id {
		t.Error("Wrong session id")
	}
	if storeSession.Value.(int) != session.Value.(int) {
		t.Error("Wrong session value")
	}
	if storeSession.Expires != session.Expires {
		t.Error("Wrong session expires date")
	}
}

func TestDelete(t *testing.T) {
	store := NewInMemorySessionStore()
	id := "abc"
	session := &Session{
		Id:      "id",
		Value:   42,
		Expires: time.Date(2013, 10, 11, 5, 20, 0, 0, time.UTC),
	}
	store.New(id, session)

	store.Delete(id)

	if len(store.sessions) != 0 {
		t.Error("Session was not deleted")
	}
}

// -----------------------------------
// Session Manager
// -----------------------------------

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
	w := newRecorder()

	manager := NewSessionManager(store, timeout)
	session := manager.initSession(w)

	header := w.Header().Get("Set-Cookie")
	if header == "" {
		t.Error("Cookie should be set")
	}
	if session == nil {
		t.Error("initSession should return session")
	}
	if session.Value != nil {
		t.Error("New session should have a nil value")
	}
}

func TestGetSessionIdFromCookie(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour
	r := newHttpRequest("GET", "localhost:8080/")
	cookieRaw := []string{"crater.SessionId=value; Expires=Fri, 11 Oct 2013 05:20:00 UTC"}
	r.Header["Cookie"] = cookieRaw

	manager := NewSessionManager(store, timeout)
	id, found := manager.getSessionIdFromCookie(r)

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
	r := newHttpRequest("GET", "localhost:8080/")

	manager := NewSessionManager(store, timeout)
	id, found := manager.getSessionIdFromCookie(r)

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
	r := newHttpRequest("GET", "localhost:8080/")
	w := newRecorder()

	manager := NewSessionManager(store, timeout)
	session := manager.GetSession(w, r)

	if session == nil {
		t.Error("Session should be initiated")
	}
}

func TestGetSession_NoSessionAndCookie_SessionShouldBeInitiated(t *testing.T) {
	store := NewInMemorySessionStore()
	timeout := time.Hour
	r := newHttpRequest("GET", "localhost:8080/")
	w := newRecorder()

	manager := NewSessionManager(store, timeout)
	session := manager.GetSession(w, r)

	if session == nil {
		t.Error("Session should be initiated")
	}
	header := w.Header().Get("Set-Cookie")
	if header == "" {
		t.Error("Cookie should be set")
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
