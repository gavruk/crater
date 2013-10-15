package session

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	rawExpiresFormat = "Fri, 01-Jan-2001 11:11:11 +0300"
	cookieName       = "crater.SessionId"
)

// SessionManager stores sessions
type SessionManager struct {
	sessions map[string]*Session
	timeout  time.Duration
	mutex    sync.RWMutex
}

// NewSessionManager creates new instance SessionManager
// NewSessionManager should be called once per application
func NewSessionManager(timeout time.Duration) *SessionManager {
	manager := new(SessionManager)
	manager.sessions = make(map[string]*Session)
	manager.timeout = timeout

	go func(manager *SessionManager) {
		for {
			now := time.Now().UTC().Unix()
			for id, session := range manager.sessions {
				if session.Expires.UTC().Unix() < now {
					delete(manager.sessions, id)
				}
			}
			time.Sleep(time.Minute)
		}
	}(manager)
	return manager
}

// Abandon terminates session by session Id
func (manager *SessionManager) Abandon(id string) {
	delete(manager.sessions, id)
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

	session, sessionFound := manager.getSessionById(sessionId)
	if !sessionFound {
		return manager.initSession(w)
	}
	return session
}

func (manager *SessionManager) getSessionIdFromCookie(r *http.Request) (id string, found bool) {
	id = ""
	c, _ := r.Cookie(cookieName)
	if c != nil {
		return c.Value, true
	}
	return "", false
}

func (manager *SessionManager) getSessionById(id string) (session *Session, found bool) {
	session, found = manager.sessions[id]
	return
}

func (manager *SessionManager) initSession(w http.ResponseWriter) (session *Session) {
	id := manager.generateId()
	session = &Session{id, nil, time.Now().UTC().Add(manager.timeout), manager}
	manager.sessions[id] = session
	cookie := &http.Cookie{
		Name:       cookieName,
		Value:      session.Id,
		Expires:    time.Now().UTC().Add(manager.timeout),
		RawExpires: time.Now().UTC().Add(manager.timeout).Format(rawExpiresFormat),
	}
	http.SetCookie(w, cookie)
	return session
}

func (manager *SessionManager) sessionExists(id string) bool {
	_, found := manager.sessions[id]
	return found
}

func (manager *SessionManager) generateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)

}
