package entity

import (
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
	defer a.addEntity(entity)
	if oldEntity, ok := a.entities[entity.GetID()]; ok {

		if oldEntity.GetState() == entity.GetState() {
			return false // No change in state
		}
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
