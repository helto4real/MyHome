package entity_test

import (
	"testing"
	"time"

	"github.com/helto4real/MyHome/core/entity"
	h "github.com/helto4real/MyHome/helpers/test"
)

func TestGetEntityList(t *testing.T) {

	el := entity.NewEntityList()

	el.SetEntity(NewFakeEntity("id", "name", "type", "state", "attributes"))
	select {
	case <-time.After(100 * time.Millisecond):
		break
	}
	entities := el.GetEntities()
	h.Equals(t, 1, len(entities))
	h.Equals(t, true, el.Close())
	// Test that close cleaned up channel ok

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
