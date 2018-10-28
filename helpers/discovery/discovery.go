package discovery

import (
	"context"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

type ZeroConfigDiscovery struct {
}

func (a ZeroConfigDiscovery) Discover(serviceCategory string) []*zeroconf.ServiceEntry {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}
	var res []*zeroconf.ServiceEntry

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			res = append(res, entry)
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, serviceCategory, "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()

	return res
}
