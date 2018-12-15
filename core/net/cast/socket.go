package cast

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/gogo/protobuf/proto"
	"github.com/ninjasphere/go-castv2/api"
)

type TlsConnection struct {
	conn     *tls.Conn
	channels []*Channel
}

func NewTlsConnection() *TlsConnection {
	return &TlsConnection{
		conn:     nil,
		channels: make([]*Channel, 0),
	}
}

func (c *TlsConnection) NewChannel(sourceId, destinationId, namespace string) *Channel {
	channel := NewChannel(c, sourceId, destinationId, namespace)
	c.channels = append(c.channels, channel)
	return channel
}

func (c *TlsConnection) Connect(ctx context.Context, host net.IP, port int) error {
	var err error
	deadline, _ := ctx.Deadline()
	dialer := &net.Dialer{
		Deadline: deadline,
	}
	c.conn, err = tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", host, port), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("Failed to connect to Chromecast: %s", err)
	}

	go c.ReceiveLoop()

	return nil
}

func (c *TlsConnection) ReceiveLoop() {
	for {
		var length uint32
		err := binary.Read(c.conn, binary.BigEndian, &length)
		if err != nil {
			log.Printf("Failed to read packet length: %s", err)
			break
		}
		if length == 0 {
			log.Println("Empty packet received")
			continue
		}

		packet := make([]byte, length)
		i, err := io.ReadFull(c.conn, packet)
		if err != nil {
			log.Printf("Failed to read packet: %s", err)
			break
		}

		if i != int(length) {
			log.Printf("Invalid packet size. Wanted: %d Read: %d", length, i)
			break
		}

		message := &api.CastMessage{}
		err = proto.Unmarshal(packet, message)
		if err != nil {
			log.Printf("Failed to unmarshal CastMessage: %s", err)
			break
		}

		log.Printf("%s ⇐ %s [%s]: %+v",
			*message.DestinationId, *message.SourceId, *message.Namespace, *message.PayloadUtf8)

		var headers PayloadHeaders
		err = json.Unmarshal([]byte(*message.PayloadUtf8), &headers)

		if err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			break
		}

		for _, channel := range c.channels {
			channel.Message(message, &headers)
		}
	}
}

func (c *TlsConnection) Send(payload interface{}, sourceId, destinationId, namespace string) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	payloadString := string(payloadJson)
	message := &api.CastMessage{
		ProtocolVersion: api.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &sourceId,
		DestinationId:   &destinationId,
		Namespace:       &namespace,
		PayloadType:     api.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadString,
	}

	proto.SetDefaults(message)

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	log.Printf("%s ⇒ %s [%s]: %s", *message.SourceId, *message.DestinationId, *message.Namespace, *message.PayloadUtf8)

	err = binary.Write(c.conn, binary.BigEndian, uint32(len(data)))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

func (c *TlsConnection) Close() error {
	// TODO: graceful shutdown
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type PayloadHeaders struct {
	Type      string `json:"type"`
	RequestId *int   `json:"requestId,omitempty"`
}

func (h *PayloadHeaders) setRequestId(id int) {
	h.RequestId = &id
}

func (h *PayloadHeaders) getRequestId() int {
	return *h.RequestId
}
