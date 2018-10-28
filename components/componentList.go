package components

import (
	"github.com/helto4real/MyHome/components/cast"
)

func GetComponents() []interface{} {
	return []interface{}{
		&cast.Cast{}}
}
