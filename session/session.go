package session

import (
	"time"
)

// -----------------
// Session
// -----------------

// Session represents current session and its data
type Session struct {
	Id      string
	Value   interface{}
	Expires time.Time
	store   SessionStore
}

func (session *Session) Abandon() {
	session.store.Delete(session.Id)
}

// ------------------------
// Session Store
// ------------------------
type SessionStore interface {
	Get(id string) *Session
	New(id string, session *Session)
	Delete(id string)
}
