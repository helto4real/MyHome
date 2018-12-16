package contracts

// IEntity represents any entity in the system
type IEntity interface {
	GetID() string
	GetName() string
	GetState() string
	GetAttributes() string
}
type IEntityList interface {
	SetEntity(IEntity) bool
}

// ILight interface implements what you can do on a light
type ILight interface {
	TurnOn()
	TurnOff()
}

// ISwitch interface implements what you can do on a switch
type ISwitch interface {
	TurnOn()
	TurnOff()
}
