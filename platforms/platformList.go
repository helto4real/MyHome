package platforms

import (
	"github.com/helto4real/MyHome/platforms/hass"
)

func GetPlatforms() []interface{} {
	return []interface{}{
		//&cast.Cast{},
		&hass.HomeAssistantPlatform{}}
}
