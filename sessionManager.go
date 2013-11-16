package crater

import (
	"net/http"
	"sync"
	"time"
)

const (
	rawExpiresFormat  = "Fri, 01-Jan-2001 11:11:11 +0300"
	sessionCookieName = "crater.SessionId"
)

// SessionManager manage session for current request.
type SessionManager struct {
	store   SessionStore
	timeout time.Duration
	mutex   sync.RWMutex
}

// NewSessionManager creates new instance of SessionManager.
func NewSessionManager(store SessionStore, timeout time.Duration) *SessionManager {
	manager := new(SessionManager)
	manager.store = store
	manager.timeout = timeout

	return manager
}

// GetSession returns session for current request.
// If there is no session, it will be created.
func (manager *SessionManager) GetSession(req *Request, res *Response) *Session {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	sessionId, cookieFound := manager.getSessionIdFromCookie(req)
	if !cookieFound {
		return manager.initSession(res, "")
	}

	session := manager.store.Get(sessionId)
	if session == nil {
		return manager.initSession(res, sessionId)
	}
	session.store = manager.store
	return session
}

// getSessionIdFromCookie read cookie with session id.
func (manager *SessionManager) getSessionIdFromCookie(req *Request) (string, bool) {
	c, err := req.Cookie(sessionCookieName)
	if c != nil && err == nil {
		return c.Value, true
	}
	return "", false
}

// initSession initialize new session.
func (manager *SessionManager) initSession(res *Response, sessionId string) *Session {
	id := sessionId
	if id == "" {
		id = GenerateId()
	}
	session := &Session{
		Id:      id,
		Value:   nil,
		Expires: time.Now().UTC().Add(manager.timeout),
	}
	if sessionId == "" {
		cookie := &http.Cookie{
			Name:    sessionCookieName,
			Value:   session.Id,
			Expires: time.Now().UTC().Add(manager.timeout),
		}
		res.SetCookie(cookie)
	}
	manager.store.New(session.Id, session)
	session.store = manager.store
	return session
}
