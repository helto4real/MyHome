package logging

import (
	"testing"

	h "github.com/helto4real/MyHome/helpers/test"
)

func TestLog(t *testing.T) {
	var (
		lastFormatString string
		lastInput        interface{}
		nrOfInterfaces   int
	)
	// Save old default log and mock them with function to test the methods
	oldLogLn := logLn
	oldLogf := logf

	logf = func(format string, v ...interface{}) {
		lastFormatString = format
		nrOfInterfaces = len(v)
	}
	logLn = func(v ...interface{}) {
		lastInput = v
	}
	defer func() {
		logLn = oldLogLn
		logf = oldLogf
	}()

	log := DefaultLogger{}

	t.Run("LogDebug single value", func(t *testing.T) {
		log.LogDebug("Hello test")
		h.Equals(t, lastInput, []interface{}{"Hello test"})

	})
	t.Run("LogDebug multiple values", func(t *testing.T) {
		log.LogDebug("Has a %s", 123)
		h.Equals(t, lastFormatString, "Has a %s\n")
		h.Equals(t, nrOfInterfaces, 1)
	})

	t.Run("LogInformation single value", func(t *testing.T) {
		log.LogInformation("Hello info")
		h.Equals(t, lastInput, []interface{}{"Hello info"})

	})
	t.Run("LogInformation multiple values", func(t *testing.T) {
		log.LogInformation("Has a %s", 123, "ksksk")
		h.Equals(t, lastFormatString, "Has a %s\n")
		h.Equals(t, nrOfInterfaces, 2)
	})

	t.Run("LogWarning single value", func(t *testing.T) {
		log.LogWarning("Hello warn")
		h.Equals(t, lastInput, []interface{}{"Hello warn"})

	})
	t.Run("LogWarning multiple values", func(t *testing.T) {
		log.LogWarning("Has a %s", 123, "ksksk", 12.2)
		h.Equals(t, lastFormatString, "Has a %s\n")
		h.Equals(t, nrOfInterfaces, 3)
	})

	t.Run("LogError single value", func(t *testing.T) {
		log.LogError("Hello error")
		h.Equals(t, lastInput, []interface{}{"Hello error"})

	})
	t.Run("LogError multiple values", func(t *testing.T) {
		log.LogError("Has a %s", 123, "ksksk", 12.2, "kk")
		h.Equals(t, lastFormatString, "Has a %s\n")
		h.Equals(t, nrOfInterfaces, 4)
	})
}
