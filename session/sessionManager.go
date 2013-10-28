package session

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	rawExpiresFormat  = "Fri, 01-Jan-2001 11:11:11 +0300"
	sessionCookieName = "crater.SessionId"
)

// -----------------------------
// In-Memory Session Store
// -----------------------------

type InMemorySessionStore struct {
	sessions map[string]*Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	store := &InMemorySessionStore{}
	store.sessions = make(map[string]*Session)

	go func(store *InMemorySessionStore) {
		for {
			now := time.Now().UTC().Unix()
			for id, session := range store.sessions {
				if session.Expires.UTC().Unix() < now {
					store.Delete(id)
				}
			}
			time.Sleep(time.Minute)
		}
	}(store)

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

// -----------------------------------
// Session Manager
// -----------------------------------

type SessionManager struct {
	store   SessionStore
	timeout time.Duration
	mutex   sync.RWMutex
}

func NewSessionManager(store SessionStore, timeout time.Duration) *SessionManager {
	manager := new(SessionManager)
	manager.store = store
	manager.timeout = timeout

	return manager
}

// GetSession returns current session
// GetSession init new session if session doesn't exist
func (manager *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) *Session {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	sessionId, cookieFound := manager.getSessionIdFromCookie(r)
	if !cookieFound {
		return manager.initSession(w)
	}

	session := manager.store.Get(sessionId)
	if session == nil {
		return manager.initSession(w)
	}
	session.store = manager.store
	return session
}

func (manager *SessionManager) getSessionIdFromCookie(r *http.Request) (id string, found bool) {
	id = ""
	c, _ := r.Cookie(sessionCookieName)
	if c != nil {
		return c.Value, true
	}
	return "", false
}

func (manager *SessionManager) initSession(w http.ResponseWriter) *Session {
	id := manager.generateId()
	session := &Session{
		Id:      id,
		Value:   nil,
		Expires: time.Now().UTC().Add(manager.timeout),
	}
	cookie := &http.Cookie{
		Name:    sessionCookieName,
		Value:   session.Id,
		Expires: time.Now().UTC().Add(manager.timeout),
	}
	http.SetCookie(w, cookie)
	manager.store.New(session.Id, session)
	session.store = manager.store
	return session
}

func (manager *SessionManager) generateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)

}
