package hass

import (
	"context"
	"encoding/json"
	"log"
	"time"

	c "github.com/helto4real/MyHome/core/contracts"
	"github.com/helto4real/MyHome/core/net"
	n "github.com/helto4real/MyHome/core/net"
)

// HomeAssistantPlatform implements integration with Home Assistant
type HomeAssistantPlatform struct {
	wsClient        *n.WsClient
	log             c.ILogger
	home            c.IMyHome
	channels        *c.Channels
	wsID            int64
	getStateId      int64
	cancelDiscovery context.CancelFunc
	context         context.Context
}

// Initialize the Home Assistant platform
func (a *HomeAssistantPlatform) Initialize(home c.IMyHome) bool {
	a.channels = home.GetChannels()
	a.log = home.GetLogger()
	a.home = home

	return true
}

func (a *HomeAssistantPlatform) InitializeDiscovery() bool {
	a.log.LogInformation("START InitializeDiscovery")
	defer a.log.LogInformation("STOP InitializeDiscovery")
	defer a.home.DoneRoutine()

	a.context, a.cancelDiscovery = context.WithCancel(context.Background())
	a.wsClient = a.connectWithReconnect()
	a.wsID = 1

	for {
		select {
		case message, mc := <-a.wsClient.ReceiverChannel:
			if !mc {
				if a.wsClient.Fatal {
					a.wsClient = a.connectWithReconnect()
					if a.wsClient == nil {
						a.log.LogInformation("Ending service discovery")
						return false
					}
				} else {
					return false
				}

			}
			var result Result
			json.Unmarshal(message, &result)
			go a.handleMessage(result)
		case <-a.context.Done():
			return false
		}

	}

}

func (a *HomeAssistantPlatform) connectWithReconnect() *net.WsClient {
	for {
		config := a.home.GetConfig()

		client := n.ConnectWS(config.HomeAssistant.IP, "/api/websocket", config.HomeAssistant.SSL)
		if client == nil {
			a.log.LogInformation("Fail to connect, reconnecting to Home Assistant in 30 seconds...")
			// Fail to connect wait to connect again
			select {
			case <-time.After(30 * time.Second):

			case <-a.context.Done():
				return nil
			}

		} else {
			return client
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

	if message.MessageType == "auth_required" {
		log.Print("message->: ", message)
		a.log.LogInformation("Got auth required, sending auth token")
		config := a.home.GetConfig()

		a.wsClient.SendString("{\"type\": \"auth\",\"access_token\": \"" + config.HomeAssistant.Token + "\"}")
	} else if message.MessageType == "auth_ok" {
		log.Print("message->: ", message)
		a.log.LogInformation("Got auth_ok, downloading all states initially")
		a.sendMessage("get_states")
	} else if message.MessageType == "result" {

		if message.Id == a.getStateId {
			a.log.LogInformation("Got all states, getting events")
			for _, data := range message.Result {
				newHassEntity := NewHassEntity("hass_"+data.EntityId, data.EntityId, "hass", data.State, data.Attributes)
				message := c.NewMessage(c.MessageType.EntityUpdated, newHassEntity)
				a.home.GetChannels().MainChannel <- *message
			}

			a.subscribeEvents()
		}
	} else if message.MessageType == "event" {
		data := message.Event.Data
		a.log.LogInformation("---------------------------------------")
		a.log.LogInformation("message->: %s", data.EntityId, data.NewState.State)
		a.log.LogInformation("---------------------------------------")
		newHassEntity := NewHassEntity("hass_"+data.EntityId, data.EntityId, "hass", data.NewState.State, data.NewState.Attributes)
		message := c.NewMessage(c.MessageType.EntityUpdated, newHassEntity)
		a.home.GetChannels().MainChannel <- *message

	}

}

func (a *HomeAssistantPlatform) EndDiscovery() {
	a.log.LogInformation("START EndDiscovery")
	defer a.log.LogInformation("STOP EndDiscovery")
	defer a.home.DoneRoutine()

	a.cancelDiscovery()
	a.wsClient.Close(false)
}
