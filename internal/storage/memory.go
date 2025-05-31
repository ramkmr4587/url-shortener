package storage

import (
	"sync"
)

type Store struct {
	URLToShort map[string]string
	ShortToURL map[string]string
	DomainHits map[string]int
	Mutex      sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		URLToShort: make(map[string]string),
		ShortToURL: make(map[string]string),
		DomainHits: make(map[string]int),
	}
}
