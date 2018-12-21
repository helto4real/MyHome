package entity_test

import (
	"testing"

	c "github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/core/entity"
	h "github.com/helto4real/MyHome/helpers/test"
)

// TestGetEntites tests the empty and one entity added
func TestGetEntities(t *testing.T) {

	el := entity.NewEntityList()

	h.Equals(t, (chan c.IEntity)(nil), el.GetEntities())

	el.SetEntity(NewFakeEntity("id", "name", "type", "state", "attributes"))

	entChannel := el.GetEntities()
	h.NotEquals(t, (chan c.IEntity)(nil), entChannel)
	h.Equals(t, 1, len(entChannel))

	ent := <-entChannel
	h.Equals(t, "id", ent.GetID())
	h.Equals(t, "name", ent.GetName())
	h.Equals(t, "type", ent.GetType())
	h.Equals(t, "state", ent.GetState())
	h.Equals(t, "attributes", ent.GetAttributes())

}

func TestGetEntitiesMultiple(t *testing.T) {
	el := entity.NewEntityList()
	el.SetEntity(NewFakeEntity("id", "name", "type", "state", "attributes"))
	h.Equals(t, 1, len(el.GetEntities()))
	el.SetEntity(NewFakeEntity("id2", "name2", "type2", "state2", "attributes2"))
	h.Equals(t, 2, len(el.GetEntities()))
}

// TestGetEntitiesSameName test that entity get replaced if same id exists
func TestGetEntitiesSameName(t *testing.T) {
	el := entity.NewEntityList()
	el.SetEntity(NewFakeEntity("id", "name", "type", "state", "attributes"))
	el.SetEntity(NewFakeEntity("id", "name2", "type2", "state2", "attributes2"))
	entChannel := el.GetEntities()
	h.Equals(t, 1, len(entChannel))
	h.Equals(t, "name2", (<-entChannel).GetName())
}

func TestHandleMessageEntityTypeCorrect(t *testing.T) {
	el := entity.NewEntityList()
	msg := c.NewMessage(c.MessageType.Entity, NewFakeEntity("id", "name", "type", "state", "attributes"))
	h.Equals(t, true, el.HandleMessage(msg))

	entChannel := el.GetEntities()
	h.NotEquals(t, (chan c.IEntity)(nil), entChannel)
	h.Equals(t, 1, len(entChannel))

	ent := <-entChannel
	h.Equals(t, "id", ent.GetID())
	h.Equals(t, "name", ent.GetName())
	h.Equals(t, "type", ent.GetType())
	h.Equals(t, "state", ent.GetState())
	h.Equals(t, "attributes", ent.GetAttributes())
}

func TestHandleMessageWrongEntityType(t *testing.T) {
	el := entity.NewEntityList()
	msg := c.NewMessage(c.MessageType.Entity, NonEntityType{})
	h.Equals(t, false, el.HandleMessage(msg))
	// We got a non entity message type so it should be nil
	h.Equals(t, (chan c.IEntity)(nil), el.GetEntities())
}

type NonEntityType struct{}

type FakeEntity struct {
	ID         string
	Name       string
	Type       string
	State      string
	Attributes string
}

// GetID returns unique id of entity
func (a FakeEntity) GetID() string         { return a.ID }
func (a FakeEntity) GetState() string      { return a.State }
func (a FakeEntity) GetType() string       { return a.Type }
func (a FakeEntity) GetAttributes() string { return a.Attributes }
func (a FakeEntity) GetName() string       { return a.Name }

func NewFakeEntity(id string, name string, entityType string, state string, attributes string) FakeEntity {
	return FakeEntity{
		ID:         id,
		Name:       name,
		Type:       entityType,
		State:      state,
		Attributes: attributes}
}
