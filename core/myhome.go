package core

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/helto4real/MyHome/core/logging"
)

type Config struct {
	Path        string
	MainChannel chan (bool)
	OsSignals   chan (os.Signal)
}

// IMyHome is the interface for main AutoHome object
type IMyHome interface {
	// Init the home automation
	Init()
	Loop()
	GetConfig() Config
}

// Init the automations
func Init(loggerUsed logging.ILogger) bool {
	logger = loggerUsed
	GetConfig()
	signal.Notify(config.OsSignals, syscall.SIGTERM)
	signal.Notify(config.OsSignals, syscall.SIGINT)
	return true
}

func Loop() bool {
	go eventHandler()

	for {
		select {
		case <-config.MainChannel:
			LogInformation("main channel")
		case s := <-config.OsSignals:
			close(config.MainChannel)
			LogInformation("SIGNAL")
			LogInformation(s.String())
			time.Sleep(2 * time.Second)
			return true
		}
	}
}

func eventHandler() {
	for {
		select {
		case _, mc := <-config.MainChannel:
			if !mc {
				LogInformation("Eventbus terminating, exiting eventhandler")
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

var logger logging.ILogger

func LogError(text string) {
	logger.LogError(text)
}

func LogWarning(text string) {
	logger.LogWarning(text)
}

func LogInformation(text string) {
	logger.LogInformation(text)
}

func LogDebug(text string) {
	logger.LogDebug(text)
}
