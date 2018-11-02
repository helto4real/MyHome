package contracts

import (
	"os"
)

type Config struct {
	Path        string
	MainChannel chan (Message)
	StopChannel chan bool
	OsSignals   chan (os.Signal)
}

type messageType struct {
	EntityUpdated string
}

var MessageType messageType = messageType{
	EntityUpdated: "entity_updated"}

type Message struct {
	Type string
	Body interface{}
}

func NewMessage(messageType string, body interface{}) *Message {
	return &Message{
		Type: messageType,
		Body: body}

}
