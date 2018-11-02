package main

import (
	"log"

	"github.com/helto4real/MyHome/core"
	"github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/core/logging"
)

func main() {

	var home contracts.IMyHome = new(core.MyHome)
	home.Init(logging.DefaultLogger{})

	home.Loop()
	log.Print("Ended")
}
