package config

import (
	"testing"
)

var yml = `
home_assistant:
    ip: '192.168.0.100'
    token: 'ABCDEFG1234567'
`
var failyml = `
home_assistant;
    ip: '192.168.0.100'
    token: ABCDEFG1234567
`

func TestOpenRawConfig(t *testing.T) {
	c, err := getRawConfig([]byte(yml))

	//Ok(t, err)

	if err != nil {
		t.Error("Failed to get raw config", err)
	}
	if c.HomeAssistant.IP != "192.168.0.100" {
		t.Error("Ip address error, expected 192.168.0.100 got ", c.HomeAssistant.IP)
	}
	if c.HomeAssistant.Token != "ABCDEFG1234567" {
		t.Error("Token error, expected ABCDEFG1234567 got ", c.HomeAssistant.Token)
	}
}

func TestMalformatedYaml(t *testing.T) {
	_, err := getRawConfig([]byte(failyml))
	if err == nil {
		t.Error("Failed failyml config should return err", err)
	}

	if err.Error() != "yaml: line 3: mapping values are not allowed in this context" {
		t.Error("Unexpected error failyaml, got", err.Error())
	}

}
