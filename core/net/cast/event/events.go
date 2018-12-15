package event

type Disconnected struct {
	Reason error
}

type Connected struct{}

type AppStarted struct {
	AppID       string
	DisplayName string
	StatusText  string
}

type AppStopped struct {
	AppID       string
	DisplayName string
	StatusText  string
}

type StatusUpdated struct {
	Level float64
	Muted bool
}
