package hass

import (
	"encoding/json"
	"log"

	c "github.com/helto4real/MyHome/core/contracts"
	n "github.com/helto4real/MyHome/core/net"
)

// HomeAssistantPlatform implements integration with Home Assistant
type HomeAssistantPlatform struct {
	wsClient   *n.WsClient
	log        c.ILogger
	home       c.IMyHome
	config     *c.Config
	wsID       int64
	getStateId int64
}

// Initialize the Home Assistant platform
func (a *HomeAssistantPlatform) Initialize(home c.IMyHome) bool {
	a.config = home.GetConfig()
	a.log = home.GetLogger()
	a.home = home

	return true
}

func (a *HomeAssistantPlatform) InitializeDiscovery() bool {
	a.log.LogInformation("START InitializeDiscovery")
	defer a.log.LogInformation("STOP InitializeDiscovery")
	defer a.home.DoneRoutine()

	a.wsClient = n.ConnectWS("192.168.1.5:8123")
	a.wsID = 1

	for {
		select {
		case message, mc := <-a.wsClient.ReceiverChannel:
			if !mc {
				a.log.LogInformation("Ending service discovery")
				return false
			}
			var result Result
			json.Unmarshal(message, &result)
			go a.handleMessage(result)
		}
	}

}

// body map[string]interface{}
func (a *HomeAssistantPlatform) sendMessage(messageType string) {
	a.wsID = a.wsID + 1
	s := map[string]interface{}{
		"id":   a.wsID,
		"type": messageType}

	if messageType == "get_states" {
		a.getStateId = a.wsID
	}
	a.wsClient.SendMap(s)

}

func (a *HomeAssistantPlatform) subscribeEvents() {
	a.wsID = a.wsID + 1
	s := map[string]interface{}{
		"id":         a.wsID,
		"type":       "subscribe_events",
		"event_type": "state_changed"}

	a.wsClient.SendMap(s)

}

func (a *HomeAssistantPlatform) handleMessage(message Result) {
	//log.Printf("message->: %s", message)

	if message.MessageType == "auth_required" {
		log.Printf("message->: %s", message)
		a.log.LogInformation("Got auth required, sending auth token")
		a.wsClient.SendString("{\"type\": \"auth\",\"access_token\": \"SomeToken\"}")
	} else if message.MessageType == "auth_ok" {
		log.Printf("message->: %s", message)
		a.log.LogInformation("Got auth_ok, downloading all states initially")
		a.sendMessage("get_states")
	} else if message.MessageType == "result" {

		if message.Id == a.getStateId {
			a.log.LogInformation("Got all states, getting events")
			a.subscribeEvents()
		}
	} else if message.MessageType == "event" {
		data := message.Event.Data
		log.Printf("\r\n\r\n---------------------------------------")
		log.Printf("message->: %s, new state: %s", data.EntityId, data.NewState.State)
		log.Printf("\r\n---------------------------------------")
	}

}

func (a *HomeAssistantPlatform) EndDiscovery() {
	a.log.LogInformation("START EndDiscovery")
	defer a.log.LogInformation("STOP EndDiscovery")
	defer a.home.DoneRoutine()
	close(a.wsClient.ReceiverChannel)
}
