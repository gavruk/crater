package crater

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

func (manager *SessionManager) GetSession(req *Request, res *Response) *Session {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	sessionId, cookieFound := manager.getSessionIdFromCookie(req)
	if !cookieFound {
		return manager.initSession(res)
	}

	session := manager.store.Get(sessionId)
	if session == nil {
		return manager.initSession(res)
	}
	session.store = manager.store
	return session
}

func (manager *SessionManager) getSessionIdFromCookie(req *Request) (id string, found bool) {
	id = ""
	c, _ := req.Cookie(sessionCookieName)
	if c != nil {
		return c.Value, true
	}
	return "", false
}

func (manager *SessionManager) initSession(res *Response) *Session {
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
	res.SetCookie(cookie)
	manager.store.New(session.Id, session)
	session.store = manager.store
	return session
}

func (manager *SessionManager) generateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)

}
