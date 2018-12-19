package logging

import "log"

// DefaultLogger logs to the standard io
type DefaultLogger struct{}

var logLn = log.Println
var logf = log.Printf

func (a DefaultLogger) LogError(format string, c ...interface{}) {
	if c == nil {
		logLn(format)
	} else {
		logf(format+"\n", c...)
	}
}

func (a DefaultLogger) LogWarning(format string, c ...interface{}) {
	if c == nil {
		logLn(format)
	} else {
		logf(format+"\n", c...)
	}
}

func (a DefaultLogger) LogInformation(format string, c ...interface{}) {
	if c == nil {
		logLn(format)
	} else {
		logf(format+"\n", c...)
	}
}

func (a DefaultLogger) LogDebug(format string, c ...interface{}) {
	if c == nil {
		logLn(format)
	} else {
		logf(format+"\n", c...)
	}
}
