package cast

import (
	"log"
	"sync"

	"github.com/grandcat/zeroconf"
	c "github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/helpers/discovery"
)

type Cast struct {
	zeroConf         *discovery.ZeroConfigDiscovery
	log              c.ILogger
	home             c.IMyHome
	config           *c.Config
	discoveryChannel chan *zeroconf.ServiceEntry
	syncRoutines     *sync.WaitGroup
}

func (a *Cast) Initialize(home c.IMyHome) bool {
	a.config = home.GetConfig()
	a.log = home.GetLogger()
	a.home = home
	a.syncRoutines = &sync.WaitGroup{}

	return true
}

func (a *Cast) doDicoveries() {
	defer a.log.LogInformation("STOP: CastDevice Discovery loop")
	defer a.syncRoutines.Done()

	a.log.LogInformation("START: CastDevice Discovery loop")
	for {
		a.zeroConf = &discovery.ZeroConfigDiscovery{}

		a.zeroConf.Discover("_googlecast._tcp", a.discoveryChannel)
		if a.zeroConf.ShouldCancel() == true {

			return
		}
	}
}

func (a *Cast) InitializeDiscovery() bool {
	log.Print("START InitializeDiscovery")
	defer log.Print("STOP InitializeDiscovery")

	defer a.home.DoneRoutine()

	a.discoveryChannel = make(chan *zeroconf.ServiceEntry)
	a.syncRoutines.Add(1)
	go a.doDicoveries()

	for {
		select {
		case entry, mc := <-a.discoveryChannel:
			if !mc {
				a.log.LogInformation("Ending service discovery")
				return false
			}
			//a.log.LogInformation("----------")
			//a.log.LogInformation("IP: %s", entry.AddrIPv4)
			//a.log.LogInformation("Host: %s", entry.HostName)
			//a.log.LogInformation("Instance: %s", entry.Instance)
			//a.log.LogInformation("Port: %s", entry.Port)
			//a.log.LogInformation("Service: %s", entry.Service)
			//a.log.LogInformation("Service I Name: %s", entry.ServiceInstanceName())
			//a.log.LogInformation("Text: %s", entry.Text)
			//a.log.LogInformation("Record: %s", entry.ServiceRecord)
			//a.log.LogInformation("Service Type: %s", entry.ServiceTypeName)
			//a.log.LogInformation("Service Type: %s", entry.ServiceTypeName)
			//a.log.LogInformation("----------")

			newCastEntity := NewCastEntity(entry.Instance, entry.AddrIPv4[0].String(), entry.Port)
			message := c.NewMessage(c.MessageType.EntityUpdated, newCastEntity)
			a.config.MainChannel <- *message

		}
	}

	return true
}

func (a *Cast) EndDiscovery() {
	log.Print("START EndDiscovery")
	defer log.Print("STOP EndDiscovery")
	defer a.home.DoneRoutine()

	a.zeroConf.EndDiscovery()
	close(a.discoveryChannel)
	(*a).syncRoutines.Wait()

}

type CastEntity struct {
	Name string
	Id   string
}

func NewCastEntity(name string, ip string, port int) *CastEntity {
	return &CastEntity{
		Name: name}
}
