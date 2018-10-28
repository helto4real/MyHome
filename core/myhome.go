package core

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/helto4real/MyHome/components"

	"github.com/helto4real/MyHome/core/contracts"
)

type Config struct {
	Path        string
	MainChannel chan (bool)
	OsSignals   chan (os.Signal)
}

type MyHome struct {
	components []interface{}
	logger     *contracts.ILogger
}

// Init the automations
func (a *MyHome) Init(loggerUsed contracts.ILogger) bool {
	a.logger = &loggerUsed
	a.components = components.GetComponents()
	a.initializeComponents()

	GetConfig()
	signal.Notify(config.OsSignals, syscall.SIGTERM)
	signal.Notify(config.OsSignals, syscall.SIGINT)
	a.setupDiscovery()
	return true
}

func (a *MyHome) Logger() contracts.ILogger {
	return *a.logger
}

func (a *MyHome) initializeComponents() {
	for _, comp := range a.components {
		x, ok := comp.(contracts.IComponent)

		if ok {
			go x.Initialize(a)
		}

	}
}
func (a *MyHome) setupDiscovery() {

	for _, comp := range a.components {
		x, ok := comp.(contracts.IDiscovery)

		if ok {
			go x.InitializeDiscovery()
		}

	}
}

func (a *MyHome) Loop() bool {
	go a.eventHandler()

	for {
		select {
		case <-config.MainChannel:
			(*a.logger).LogInformation("main channel")
		case s := <-config.OsSignals:
			close(config.MainChannel)
			(*a.logger).LogInformation("SIGNAL")
			(*a.logger).LogInformation(s.String())
			time.Sleep(2 * time.Second)
			return true
		}
	}
}

func (a *MyHome) eventHandler() {
	for {
		select {
		case _, mc := <-config.MainChannel:
			if !mc {
				(*a.logger).LogInformation("Eventbus terminating, exiting eventhandler")
				return
			}
		case <-time.After(1 * time.Second):
			config.MainChannel <- true
		}
	}
}

var config *Config

func GetConfig() Config {
	if config == nil {
		config = &Config{
			"Hello",
			make(chan bool),
			make(chan os.Signal, 1)}
	}
	return *config
}
