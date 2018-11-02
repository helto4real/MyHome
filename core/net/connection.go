package net

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/helto4real/MyHome/core/contracts"
)

var addr = flag.String("addr", ":8082", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "./front-end/build/es5-bundled/index.html")
}

func listenEvents(hub *Hub, config *contracts.Config) {
	log.Print("START: Listen to events")
	defer log.Print("STOP: Listen to events")
	for {
		select {
		case msg, ok := <-config.MainChannel:
			if !ok {
				log.Print("Main channel terminating, exiting ListenEvents")
				return
			}
			jsonMessage, err := json.Marshal(msg)
			if err != nil {
				log.Print("Failed to convert message to json", err)
			}

			//log.Print("Event:", string(jsonMessage[:]))
			hub.SendMessage(jsonMessage)
			// case <-config.OsSignals:
			// 	log.Print("Closing websocket sender")
			// 	return
		}
	}
}

func CloseWebServers() {
	hub.CloseHub()
	cancelWebServer()
}

var hub *Hub
var cancelWebServer context.CancelFunc

func SetupWebservers(home contracts.IMyHome) {
	defer home.DoneRoutine()

	log.Print("START: Webserver")
	defer log.Print("STOP: Webserver")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cancelWebServer = cancel

	webServer := &http.Server{Addr: ":8082"}
	hub = NewHub()
	go hub.Run()
	go listenEvents(hub, home.GetConfig())

	http.Handle("/", http.FileServer(http.Dir("./front-end/build/es5-bundled")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	go func() {
		err := webServer.ListenAndServe()
		if err != nil {
			log.Print("ListenAndServe: ", err)
		}
	}()
	select {
	case <-ctx.Done():
		// Shutdown the server when the context is canceled
		err := webServer.Shutdown(ctx)
		if err != nil {
			log.Print(err)
		}
	}

}
