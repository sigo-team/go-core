package lib

type IdentifierManager struct {
	lastInsertedID int64
}

func NewIdentifierManager() *IdentifierManager {
	return &IdentifierManager{
		lastInsertedID: 10000,
	}
}

func (m *IdentifierManager) NewID() int64 {
	defer func() {
		m.lastInsertedID++
	}()
	return m.lastInsertedID
}
