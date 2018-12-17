package config_test

import (
	"testing"

	c "github.com/helto4real/MyHome/core/config"
	h "github.com/helto4real/MyHome/helpers/test"
)

type mockos struct {
}

func (o mockos) HomePath() string {

	return "testdata"
}

type mockosFail struct {
}

func (o mockosFail) HomePath() string {

	return ""
}

func TestOpen(t *testing.T) {
	configuration := c.NewConfigurationMock(mockos{})
	config, _ := configuration.Open()
	h.Equals(t, config.HomeAssistant.IP, "192.168.0.100")
}

func TestFailOpenConfigFile(t *testing.T) {
	configuration := c.NewConfigurationMock(mockosFail{})
	config, err := configuration.Open()
	h.Assert(t, config == nil, "Configuration should return nil value")
	h.Assert(t, err != nil, "Configuration should return error value")
}
func TestNewConfiguration(t *testing.T) {
	h.Assert(t, c.NewConfiguration() != nil, "Configuration failed")
}
