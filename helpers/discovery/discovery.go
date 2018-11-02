package discovery

import (
	"context"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

type ZeroConfigDiscovery struct {
	cancel       context.CancelFunc
	shouldCancel bool
}

func (a *ZeroConfigDiscovery) EndDiscovery() {
	a.shouldCancel = true
	a.cancel()
}

func (a *ZeroConfigDiscovery) ShouldCancel() bool {
	return a.shouldCancel
}

func (a *ZeroConfigDiscovery) Discover(serviceCategory string, discoveredDevices chan *zeroconf.ServiceEntry) {
	log.Print("START Discovery Zero Conf")
	defer log.Print("END Discovery Zero Conf")

	a.shouldCancel = false
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			discoveredDevices <- entry
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	a.cancel = cancel
	defer cancel()
	err = resolver.Browse(ctx, serviceCategory, "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
}
