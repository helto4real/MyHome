package contracts

// IComponent represents any component in the system
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
	GetLogger() ILogger
	GetEntityList() IEntityList
	GetConfig() *Config
	StartRoutine()
	DoneRoutine()
}

type ILogger interface {
	// Init the home automation
	LogError(format string, a ...interface{})
	LogWarning(format string, a ...interface{})
	LogInformation(format string, a ...interface{})
	LogDebug(format string, a ...interface{})
}
