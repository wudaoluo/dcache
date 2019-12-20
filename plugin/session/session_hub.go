package session

import (
	"sync"
)
type SessHub interface {
	Add(sess *Session)
	Del(ID string)
}

type sessionHub struct {
	mu sync.RWMutex
	hub map[string]*Session
}

func NewSessionHub() SessHub{
	return &sessionHub{hub:make(map[string]*Session)}
}

func (s *sessionHub) Add(sess *Session) {
	s.mu.Lock()
	s.hub[sess.ID()] = sess
	s.mu.Unlock()
}

func (s *sessionHub) Del(ID string) {
	s.mu.Lock()
	delete(s.hub,ID)
	s.mu.Unlock()
}