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
)

type SessionManager struct {
	sessions map[string]*Session
	timeout  time.Duration
	mutex    sync.RWMutex
}

func NewSessionManager(timeout time.Duration) *SessionManager {
	manager := new(SessionManager)
	manager.sessions = make(map[string]*Session)
	manager.timeout = timeout

	go func(manager *SessionManager) {
		for {
			now := time.Now().UTC().Unix()
			for id, session := range manager.sessions {
				if session.expires.UTC().Unix() < now {
					delete(manager.sessions, id)
				}
			}
			time.Sleep(time.Minute)
		}
	}(manager)
	return manager
}

func (manager *SessionManager) Abandon(id string) {
	delete(manager.sessions, id)
}

func (manager *SessionManager) getSessionById(id string) (session *Session) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	if id == "" || !manager.sessionExists(id) {
		id = manager.generateId()
	}
	s, found := manager.sessions[id]
	if !found {
		session = &Session{id, nil, time.Now().UTC().Add(manager.timeout), manager}
		manager.sessions[id] = session
	} else {
		session = s
	}
	return
}

func (manager *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) (session *Session) {
	sessionId := ""
	c, _ := r.Cookie("crater.SessionId")
	if c != nil {
		sessionId = c.Value
	}

	session = manager.getSessionById(sessionId)

	cookie := &http.Cookie{
		Name:       "crater.SessionId",
		Value:      session.Id,
		Expires:    time.Now().Add(time.Hour),
		RawExpires: time.Now().Add(time.Hour).Format(rawExpiresFormat),
	}
	http.SetCookie(w, cookie)
	return
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
