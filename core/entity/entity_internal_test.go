package entity

import (
	"testing"

	"github.com/helto4real/MyHome/core/contracts"
	h "github.com/helto4real/MyHome/helpers/test"
)

func TestConstructorAndClose(t *testing.T) {

	el := NewEntityList()
	h.NotEquals(t, nil, el.enitityChannel)
	h.NotEquals(t, nil, el.entities)
	h.Equals(t, 0, len(el.entities))

	h.Equals(t, true, el.Close())
	// Test that close cleaned up channel ok
	h.Equals(t, map[string]contracts.IEntity(nil), el.entities)

	h.Equals(t, false, el.Close())
}

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

func NewFakeEntity(id string, name string, entityType string, state string, attributes string) *FakeEntity {
	return &FakeEntity{
		ID:         id,
		Name:       name,
		Type:       entityType,
		State:      state,
		Attributes: attributes}
}
