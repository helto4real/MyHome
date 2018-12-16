package platforms

import (
	"github.com/helto4real/MyHome/platforms/cast"
)

func GetPlatforms() []interface{} {
	return []interface{}{
		&cast.Cast{}}
}
