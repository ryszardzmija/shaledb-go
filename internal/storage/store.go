package storage

import (
	"bytes"
	"errors"
	"sync"
)

type Store interface {
	Get(key []byte) (value []byte, found bool, err error)
	Put(key, value []byte) (err error)
	Delete(key []byte) (found bool, err error)
}

type InMemoryStore struct {
	mu      sync.RWMutex
	hashMap map[string][]byte
}

func (s *InMemoryStore) Get(key []byte) (value []byte, found bool, err error) {
	if len(key) == 0 {
		return nil, false, errors.New("key must not be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.hashMap[string(key)]
	if !ok {
		return nil, false, nil
	}

	return bytes.Clone(v), true, nil
}

func (s *InMemoryStore) Put(key, value []byte) (err error) {
	if len(key) == 0 {
		return errors.New("key must not be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.hashMap == nil {
		s.hashMap = make(map[string][]byte)
	}

	s.hashMap[string(key)] = bytes.Clone(value)

	return nil
}

func (s *InMemoryStore) Delete(key []byte) (found bool, err error) {
	if len(key) == 0 {
		return false, errors.New("key must not be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	keyStr := string(key)

	_, ok := s.hashMap[keyStr]
	if !ok {
		return false, nil
	}

	delete(s.hashMap, keyStr)

	return true, nil
}
