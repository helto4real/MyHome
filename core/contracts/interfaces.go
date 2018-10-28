package contracts

// IDevice represents any device in the system
type IComponent interface {
	Initialize(home IMyHome) bool
}

// IDiscovery represents any device in the system
type IDiscovery interface {
	InitializeDiscovery() bool
	EndDiscovery()
}

// IMyHome is the interface for main AutoHome object
type IMyHome interface {
	// Init the home automation
	Init(ILogger) bool
	Loop() bool
	Logger() ILogger
	//GetConfig() Config
}

// IDevice represents any device in the system
type IDevice interface {
	Id() string
	Name() string
}

type ILogger interface {
	// Init the home automation
	LogError(format string, a ...interface{})
	LogWarning(format string, a ...interface{})
	LogInformation(format string, a ...interface{})
	LogDebug(format string, a ...interface{})
}
