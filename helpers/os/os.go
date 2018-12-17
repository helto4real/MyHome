package os

import (
	"log"
	"os/user"

	"github.com/helto4real/MyHome/core/contracts"
)

type osHelper struct{}

func (o osHelper) HomePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Print("Failed to get user!")
		return ""
	}
	return usr.HomeDir
}

func NewOsHelper() contracts.IOS {
	return osHelper{}
}
