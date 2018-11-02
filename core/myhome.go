package core

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/helto4real/MyHome/components"

	"github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/core/net"
)

type MyHome struct {
	components   []interface{}
	logger       contracts.ILogger
	syncRoutines sync.WaitGroup
}

// Init the automations
func (a *MyHome) Init(loggerUsed contracts.ILogger) bool {
	a.syncRoutines = sync.WaitGroup{}
	a.logger = loggerUsed
	newConfig()
	a.components = components.GetComponents()
	a.initializeComponents()

	signal.Notify(config.OsSignals, syscall.SIGTERM)
	signal.Notify(config.OsSignals, syscall.SIGINT)
	a.setupWebservers()
	a.setupDiscovery()
	return true
}

func (a *MyHome) end() {
	a.endDiscovery()
	net.CloseWebServers()
	// Wait for the main GoRoutines to finish
	a.logger.LogInformation("Wait for the main GoRoutines to finish")
	if a.waitToEnd() {
		a.logger.LogInformation("Not all goroutines closed, forcing end.")

	} else {
		a.logger.LogInformation("All goroutines ended, closing application")
	}
	// Wait some additional time to see debug messages on go routine shutdown.
	//time.Sleep(5 * time.Second)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func (a *MyHome) waitToEnd() bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		a.waitRoutines()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(time.Second * 10):
		return true // timed out
	}
}

var count int

func (a *MyHome) StartRoutine() {
	count = count + 1

	a.syncRoutines.Add(1)
	log.Print("Counter", count)

}
func (a *MyHome) DoneRoutine() {
	count = count - 1
	a.syncRoutines.Done()
	log.Print("Counter", count)
}
func (a *MyHome) waitRoutines() {
	log.Print("Wait Counter", count)
	a.syncRoutines.Wait()
}

func (a *MyHome) GetLogger() contracts.ILogger {
	return a.logger
}

func (a *MyHome) initializeComponents() {
	for _, comp := range a.components {
		x, ok := comp.(contracts.IComponent)

		if ok {
			var h contracts.IMyHome = a
			x.Initialize(h)
		}

	}
}
func (a *MyHome) setupDiscovery() {

	for _, comp := range a.components {
		x, ok := comp.(contracts.IDiscovery)

		if ok {
			a.StartRoutine()
			go x.InitializeDiscovery()
		}

	}
}

func (a *MyHome) endDiscovery() {
	for _, comp := range a.components {
		x, ok := comp.(contracts.IDiscovery)

		if ok {
			a.StartRoutine()
			go x.EndDiscovery()
		}

	}
}

func (a *MyHome) setupWebservers() {
	a.StartRoutine()
	go net.SetupWebservers(a)
}

func (a *MyHome) Loop() bool {
	a.logger.LogInformation("Starting main LOOP")
	defer a.logger.LogInformation("Ending main LOOP")

	a.StartRoutine()
	go a.eventHandler()

	for {
		select {
		case _, mc := <-config.MainChannel:
			if !mc {
				a.logger.LogInformation("Main channel terminating, exiting Loop")
				return false
			}
		case <-config.OsSignals:
			a.logger.LogInformation("OS SIGNAL")
			close(config.MainChannel)
			a.end()

			return true
		case <-config.StopChannel:
			close(config.MainChannel)
			a.end()

			return true
		}
	}
}

func (a *MyHome) eventHandler() {
	defer a.DoneRoutine()

	for {
		select {
		case _, mc := <-config.MainChannel:
			if !mc {
				a.logger.LogInformation("Eventbus terminating, exiting eventhandler")
				return
			}
			// case <-time.After(5 * time.Second):
			// 	config.StopChannel <- true

		}
	}
}

func (a *MyHome) GetConfig() *contracts.Config {
	return config
}

var config *contracts.Config

func newConfig() contracts.Config {
	if config == nil {
		config = &contracts.Config{
			Path:        "Hello",
			MainChannel: make(chan contracts.Message),
			StopChannel: make(chan bool),
			OsSignals:   make(chan os.Signal, 1)}
	}
	return *config
}
