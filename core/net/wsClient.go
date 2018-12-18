package net

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 131072
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type WsClient struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	SendChannel     chan []byte
	ReceiverChannel chan []byte
	Fatal           bool
	syncRoutines    sync.WaitGroup
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *WsClient) readPump() {
	c.syncRoutines.Add(1)
	defer func() {
		c.syncRoutines.Done()
		c.Close(true)
		log.Printf("Close ws readpump")

	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			} else {
				log.Printf("Error reading websocket: %v", err)
			}

			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		c.ReceiverChannel <- message

	}
}

// Close the web socket client
func (c *WsClient) Close(fatal bool) {
	c.Fatal = fatal
	if c.ReceiverChannel == nil || c.SendChannel == nil {
		return
	}
	// Close the connection and ignore errors
	c.conn.Close()

	close(c.ReceiverChannel)
	c.ReceiverChannel = nil
	close(c.SendChannel)
	c.SendChannel = nil

	//  Wait for the routines to stop
	c.syncRoutines.Wait()
	log.Printf("Closing websocket to Home Assistant")
}

func (c *WsClient) SendMap(message map[string]interface{}) {

	jsonString, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshal message: %s", err)
		return
	}

	c.SendChannel <- jsonString

}

func (c *WsClient) SendString(message string) {

	c.SendChannel <- []byte(message)

}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WsClient) writePump() {
	c.syncRoutines.Add(1)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close(true)
		c.syncRoutines.Done()
	}()
	for {
		select {
		case message, ok := <-c.SendChannel:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			//log.Print("WROTE TO SEND QUEUE: ", string(message))
			// Add queued chat messages to the current websocket message.
			n := len(c.SendChannel)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.SendChannel)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ConnectWS connects to Web Sockets
func ConnectWS(ip string, path string, ssl bool) *WsClient {
	var scheme string = "ws"
	if ssl == true {
		scheme = "wss"
	}
	u := url.URL{Scheme: scheme, Host: ip, Path: path}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Print("dial:", err)
		return nil
	}

	client := &WsClient{conn: c, SendChannel: make(chan []byte, 256), ReceiverChannel: make(chan []byte), Fatal: false}

	// Do write and read operations in own go routines
	go client.writePump()
	go client.readPump()

	return client
}
