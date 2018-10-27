package logging

import (
	"log"
)

type ILogger interface {
	// Init the home automation
	LogError(text string)
	LogWarning(text string)
	LogInformation(test string)
	LogDebug(text string)
}

// DefaultLogger logs to the standard io
type DefaultLogger struct{}

func (a DefaultLogger) LogError(text string) {
	log.Fatalln(text)
}

func (a DefaultLogger) LogWarning(text string) {
	log.Println(text)
}

func (a DefaultLogger) LogInformation(text string) {
	log.Println(text)
}

func (a DefaultLogger) LogDebug(text string) {
	log.Println(text)
}
