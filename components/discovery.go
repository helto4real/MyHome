package components

import (
	"context"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

// IDiscovery represents any device in the system
type IDiscovery interface {
	InitializeDiscovery()
	EndDiscovery()
}

type ZeroConfigDiscovery struct {
}

func (a ZeroConfigDiscovery) InitializeDiscovery() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println(entry)
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, "_googlecast._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}
func (a ZeroConfigDiscovery) EndDiscovery() {

}
