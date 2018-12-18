package hass

type HassEntity struct {
	ID         string
	Name       string
	Type       string
	State      string
	Attributes string
}

// GetID returns unique id of entity
func (a HassEntity) GetID() string         { return a.ID }
func (a HassEntity) GetState() string      { return a.State }
func (a HassEntity) GetType() string       { return a.Type }
func (a HassEntity) GetAttributes() string { return a.Attributes }
func (a HassEntity) GetName() string       { return a.Name }

func NewHassEntity(id string, name string, entityType string, state string, attributes string) *HassEntity {
	return &HassEntity{
		ID:         id,
		Name:       name,
		Type:       entityType,
		State:      state,
		Attributes: attributes}
}
