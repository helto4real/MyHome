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
	entities       map[string]c.IEntity
	enitityChannel chan c.IEntity
	syncRoutines   sync.WaitGroup
	home           c.IMyHome
}

// NewEntityList makes a new instance of entity list
func NewEntityList() *List {
	el := List{
		entities:       make(map[string]c.IEntity),
		enitityChannel: make(chan c.IEntity)}

	// All edits go through own go routine
	go el.handleEntities()
	return &el
}

// GetEntities returns all entities in list
func (a *List) GetEntities() map[string]c.IEntity {
	return a.entities
}

// SetEntity returns true if not exist or state changed
func (a *List) setEntity(entity c.IEntity) bool {
	defer a.addEntity(entity)
	if oldEntity, ok := a.entities[entity.GetID()]; ok {

		if oldEntity.GetState() == entity.GetState() {
			return false // No change in state
		}
	}
	return true
}

// SetEntity returns true if not exist or state changed
func (a *List) SetEntity(entity c.IEntity) {
	a.enitityChannel <- entity
}

func (a *List) addEntity(entity c.IEntity) {
	a.entities[entity.GetID()] = entity
}

// Close the entity list go routines
func (a *List) Close() bool {
	if a.entities == nil {
		return false
	}

	close(a.enitityChannel)
	a.syncRoutines.Wait()
	a.entities = nil
	return true
}

// HandleEntities make sure all edits to the entity list
// is syncronized through own goroutine so only one thread
// can write to the list
func (a *List) handleEntities() {
	a.syncRoutines.Add(1)
	defer a.syncRoutines.Done()

	for {
		select {
		case entity, ok := <-a.enitityChannel:
			if !ok {
				return
			}
			a.setEntity(entity)
		}
	}

}

// HandleMessage handle messages from the main channel
func (a *List) HandleMessage(message c.Message) bool {
	entity, ok := message.Body.(c.IEntity)
	if !ok {
		return false
	}

	switch mt := message.Type; mt {
	case c.MessageType.Entity:
		a.enitityChannel <- entity
	}

	return true
}
