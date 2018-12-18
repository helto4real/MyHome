package cast

import (
	"log"
	"strings"
	"sync"

	"github.com/grandcat/zeroconf"
	c "github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/helpers/discovery"
)

type Cast struct {
	zeroConf         *discovery.ZeroConfigDiscovery
	log              c.ILogger
	home             c.IMyHome
	config           *c.Channels
	discoveryChannel chan *zeroconf.ServiceEntry
	syncRoutines     *sync.WaitGroup
}

func (a *Cast) Initialize(home c.IMyHome) bool {
	a.config = home.GetChannels()
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
	a.log.LogInformation("START InitializeDiscovery")
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
			a.log.LogInformation("IP: %s", entry.AddrIPv4)
			//a.log.LogInformation("Host: %s", entry.HostName)
			//a.log.LogInformation("Instance: %s", entry.Instance)
			a.log.LogInformation("Port: %s", entry.Port)
			//a.log.LogInformation("Service: %s", entry.Service)
			//a.log.LogInformation("Service I Name: %s", entry.ServiceInstanceName())
			//a.log.LogInformation("Text: %s", entry.Text)
			//a.log.LogInformation("Record: %s", entry.ServiceRecord)
			//a.log.LogInformation("Service Type: %s", entry.ServiceTypeName)
			//a.log.LogInformation("Service Type: %s", entry.ServiceTypeName)
			//a.log.LogInformation("----------")
			var deviceName = "default"
			for _, name := range entry.Text {
				s := strings.Split(name, "=")
				attribute := s[0]
				text := s[1]
				if attribute == "fn" {
					deviceName = text
					a.log.LogInformation("Found device: %s", deviceName)
				}
			}
			newCastEntity := NewCastEntity("cast_"+entry.Instance, deviceName, entry.AddrIPv4[0].String(), entry.Port)
			message := c.NewMessage(c.MessageType.EntityUpdated, newCastEntity)
			a.config.MainChannel <- *message

		}
	}

}

func (a *Cast) EndDiscovery() {
	a.log.LogInformation("START EndDiscovery")
	defer a.log.LogInformation("STOP EndDiscovery")
	defer a.home.DoneRoutine()

	a.zeroConf.EndDiscovery()
	close(a.discoveryChannel)
	(*a).syncRoutines.Wait()

}

type CastEntity struct {
	ID         string
	Name       string
	State      string
	Attributes string
	IP         string
}

// GetId returns unique id of entity
func (a *CastEntity) GetID() string         { return a.ID }
func (a *CastEntity) GetState() string      { return a.State }
func (a *CastEntity) GetType() string       { return "Cast" }
func (a *CastEntity) GetAttributes() string { return a.Attributes }
func (a *CastEntity) GetName() string       { return a.Name }

func NewCastEntity(id string, name string, ip string, port int) *CastEntity {
	return &CastEntity{
		ID:   id,
		IP:   ip,
		Name: name}
}
