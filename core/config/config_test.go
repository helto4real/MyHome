package config_test

import (
	"errors"
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
	h.Equals(t, config.HomeAssistant.SSL, false)
	h.Equals(t, config.HomeAssistant.Token, "ABCDEFG1234567")

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

type failReader struct{}

func (a failReader) Read(read []byte) (int, error) {
	return 0, errors.New("Fake error")
}
func TestOpenReaderFails(t *testing.T) {
	config := c.NewConfigurationMock(mockosFail{})
	_, err := config.OpenReader(failReader{})

	h.Assert(t, err != nil, "Expected error!")

}
