package contracts

import (
	"os"
)

type Channels struct {
	MainChannel      chan (Message)
	BroadCastChannel chan (Message)
	EventChannel     chan (Message)
	StopChannel      chan bool
	OsSignals        chan (os.Signal)
}

func (a *Channels) CloseChannels() {
	close(a.BroadCastChannel)
	close(a.EventChannel)
	close(a.StopChannel)
	close(a.MainChannel)
}

type messageType struct {
	Entity string
}

var MessageType messageType = messageType{
	Entity: "entity"}

type Message struct {
	Type string
	Body interface{}
}

func NewMessage(messageType string, body interface{}) *Message {
	return &Message{
		Type: messageType,
		Body: body}

}
