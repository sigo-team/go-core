package lib

import "sync"

type IdentifierManager struct {
	lastInsertedID int64
	mtx            sync.Mutex
}

func NewIdentifierManager() *IdentifierManager {
	return &IdentifierManager{
		lastInsertedID: 10000,
		mtx:            sync.Mutex{},
	}
}

func (m *IdentifierManager) NewID() int64 {
	m.mtx.Lock()
	defer func() {
		m.lastInsertedID++
		m.mtx.Unlock()
	}()
	return m.lastInsertedID
}
