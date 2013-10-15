package session

import (
	"time"
)

// Session represents current session and its data
type Session struct {
	Id      string
	Value   interface{}
	Expires time.Time
	manager *SessionManager
}

// Abandon terminates current session
func (session *Session) Abandon() {
	session.manager.Abandon(session.Id)
}
