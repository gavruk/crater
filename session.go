package crater

import (
	"time"
)

// Session represents current session and its data.
type Session struct {
	Id      string
	Value   interface{}
	Expires time.Time
	store   SessionStore
}

// Abandon end current session
func (session *Session) Abandon() {
	session.store.Delete(session.Id)
}

// SessionStore is an interface which allow to implement custom session store.
type SessionStore interface {
	Get(id string) *Session
	New(id string, session *Session)
	Delete(id string)
}
