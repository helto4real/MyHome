package entity

import (
	"log"

	c "github.com/helto4real/MyHome/core/contracts"
)

type EntityList struct {
	entities map[string]c.IEntity
	home     c.IMyHome
}

func NewEntityList(home c.IMyHome) EntityList {
	return EntityList{
		entities: make(map[string]c.IEntity),
		home:     home}
}

// SetEntity returns true if not exist or state changed
func (a *EntityList) SetEntity(entity c.IEntity) bool {

	if oldEntity, ok := a.entities[entity.GetID()]; ok {
		defer a.addEntity(entity)
		if oldEntity.GetState() == entity.GetState() {
			log.Printf("SAME OLD:NEW %s : %s", oldEntity.GetName(), entity.GetName())
			return false // No change in state
		}
	} else {
		a.addEntity(entity)
	}
	return true
}
func (a *EntityList) addEntity(entity c.IEntity) {
	a.entities[entity.GetID()] = entity
}

func (a *EntityList) HandleMessage(message c.Message) bool {
	entity, ok := message.Body.(c.IEntity)
	if !ok {
		return false // Todo: Add log here later
	}

	switch mt := message.Type; mt {
	case c.MessageType.EntityUpdated:
		return a.SetEntity(entity)
	}

	return true
}