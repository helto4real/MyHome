package components

// IDevice represents any device in the system
type IDevice interface {
	Id() string
	Name() string
}
