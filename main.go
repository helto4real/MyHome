package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/helto4real/MyHome/core"
	"github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/core/logging"
)

var addr = flag.String("addr", ":8082", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "front-end/build/esm-bundled/index.html")
}

func setupWebserver() {
	//http.HandleFunc("/", serveHome)
	http.Handle("/", http.FileServer(http.Dir("./front-end/build/es6-bundled")))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func main() {
	go setupWebserver()
	var home contracts.IMyHome = new(core.MyHome)
	home.Init(logging.DefaultLogger{})

	home.Loop()
}
