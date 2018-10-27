package main

import (
	"github.com/helto4real/MyHome/core"
	"github.com/helto4real/MyHome/core/logging"
)

func main() {

	core.Init(logging.DefaultLogger{})
	core.Loop()
}
