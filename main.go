package main

import (
	"log"

	"github.com/helto4real/MyHome/core"
	"github.com/helto4real/MyHome/core/config"
	"github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/core/logging"
)

/*
	"github.com/helto4real/MyHome/core"
	"github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/core/logging"
*/

func main() {
	// ip := net.ParseIP("192.168.1.133")
	// var client = n.NewClient(ip, 8009)
	// var ctx, cancel = context.WithTimeout(context.Background(), time.Second*135)

	// defer cancel()
	// defer client.Close()

	// client.Connect(ctx)

	// select {
	// case <-ctx.Done():
	// 	fmt.Println("overslept")
	// case <-time.After(130 * time.Second):
	// 	fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	// }

	var home contracts.IMyHome = new(core.MyHome)
	config, err := config.NewConfiguration().Open()
	if err != nil {
		log.Print("Failed to open config.\r\n")
		log.Print(err)
		return
	}
	home.Init(logging.DefaultLogger{}, config)

	home.Loop()
	// OsSignals := make(chan os.Signal, 1)
	// wsclient := net.ConnectWS()

	// for {
	// 	select {
	// 	case <-OsSignals:
	// 		return
	// 	}
	// }
	log.Print("ENDED")
}
