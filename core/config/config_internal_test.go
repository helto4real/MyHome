package config

import (
	"testing"

	h "github.com/helto4real/MyHome/helpers/test"
)

var yml = `
home_assistant:
    ip: '192.168.0.100'
    ssl: true
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
	h.Equals(t, err, nil)
	h.Equals(t, c.HomeAssistant.IP, "192.168.0.100")
	h.Equals(t, c.HomeAssistant.SSL, true)
	h.Equals(t, c.HomeAssistant.Token, "ABCDEFG1234567")

}

func TestMalformatedYaml(t *testing.T) {
	_, err := getRawConfig([]byte(failyml))
	h.NotEquals(t, err, nil)
	h.Equals(t, err.Error(), "yaml: line 3: mapping values are not allowed in this context")

}
