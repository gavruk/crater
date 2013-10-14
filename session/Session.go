package session

import (
	"time"
)

type Session struct {
	Id      string
	Value   interface{}
	expires time.Time
	manager *SessionManager
}

func (session *Session) Abandon() {
	session.manager.Abandon(session.Id)
}
