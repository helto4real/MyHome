package cast

import (
	"github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/helpers/discovery"
)

type Cast struct {
	zeroConf discovery.ZeroConfigDiscovery
	log      contracts.ILogger
	home     contracts.IMyHome
}

func (a *Cast) Initialize(home contracts.IMyHome) bool {
	a.log = home.Logger()
	a.home = home
	return true
}

func (a *Cast) InitializeDiscovery() bool {

	entries := a.zeroConf.Discover("_googlecast._tcp")
	for _, entry := range entries {
		a.log.LogInformation("----------")
		a.log.LogInformation("IP: %s", entry.AddrIPv4)
		a.log.LogInformation("Host: %s", entry.HostName)
		a.log.LogInformation("Instance: %s", entry.Instance)
		a.log.LogInformation("Port: %s", entry.Port)
		a.log.LogInformation("Service: %s", entry.Service)
		a.log.LogInformation("Service I Name: %s", entry.ServiceInstanceName())
		a.log.LogInformation("Text: %s", entry.Text)
		a.log.LogInformation("Record: %s", entry.ServiceRecord)
		a.log.LogInformation("Service Type: %s", entry.ServiceTypeName)
		a.log.LogInformation("Service Type: %s", entry.ServiceTypeName)
		a.log.LogInformation("----------")
	}
	return true
}

func (a *Cast) EndDiscovery() {

}
