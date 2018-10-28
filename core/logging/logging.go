package logging

import (
	"log"
)

// DefaultLogger logs to the standard io
type DefaultLogger struct{}

func (a DefaultLogger) LogError(format string, c ...interface{}) {
	if c == nil {
		log.Println()
	} else {
		log.Printf(format+"\n", c)
	}
}

func (a DefaultLogger) LogWarning(format string, c ...interface{}) {
	if c == nil {
		log.Println()
	} else {
		log.Printf(format+"\n", c)
	}
}

func (a DefaultLogger) LogInformation(format string, c ...interface{}) {
	if c == nil {
		log.Println()
	} else {
		log.Printf(format+"\n", c)
	}
}

func (a DefaultLogger) LogDebug(format string, c ...interface{}) {
	if c == nil {
		log.Println()
	} else {
		log.Printf(format+"\n", c)
	}
}
