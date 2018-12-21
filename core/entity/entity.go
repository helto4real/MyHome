package entity

import (
	"sync"

	c "github.com/helto4real/MyHome/core/contracts"
)

// List stores all entites and its states in memeory
//
// It support threadsafety by handling all writes through
// A go routine
type List struct {
	entities map[string]c.IEntity
	m        sync.Mutex
}

// NewEntityList makes a new instance of entity list
func NewEntityList() List {
	return List{entities: make(map[string]c.IEntity)}
}

// GetEntities returns a thread safe way to get all entities through a channel
//
func (a *List) GetEntities() chan c.IEntity {
	a.m.Lock()

	defer a.m.Unlock()

	if len(a.entities) == 0 {
		return nil
	}

	entityChannel := make(chan c.IEntity, len(a.entities))
	defer close(entityChannel)

	for _, entity := range a.entities {
		entityChannel <- entity
	}
	return entityChannel
}

// SetEntity returns true if not exist or state changed
func (a *List) SetEntity(entity c.IEntity) {
	a.m.Lock()
	defer a.m.Unlock()
	a.entities[entity.GetID()] = entity
}

// HandleMessage handle messages from the main channel
func (a *List) HandleMessage(message *c.Message) bool {
	entity, ok := message.Body.(c.IEntity)
	if !ok {
		return false
	}

	a.SetEntity(entity)

	return true
}
