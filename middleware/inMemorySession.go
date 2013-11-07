package middleware

import (
	"time"

	"github.com/gavruk/crater"
)

var sessionManager = crater.NewSessionManager(NewInMemorySessionStore(), time.Hour*12)

func InMemorySession(req *crater.Request, res *crater.Response) {
	req.Session = sessionManager.GetSession(req, res)
}

// -----------------------------
// In-Memory Session Store
// -----------------------------

type InMemorySessionStore struct {
	sessions map[string]*crater.Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	store := &InMemorySessionStore{}
	store.sessions = make(map[string]*crater.Session)

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

func (store InMemorySessionStore) Get(id string) *crater.Session {
	return store.sessions[id]
}

func (store InMemorySessionStore) New(id string, session *crater.Session) {
	store.sessions[id] = session
}
