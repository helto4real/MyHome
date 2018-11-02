// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"context"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	cancel context.CancelFunc
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool)}
}

func (h *Hub) SendMessage(message []byte) {
	h.broadcast <- message
}
func (h *Hub) CloseHub() {
	log.Printf("Closing hub")
	h.cancel()
}
func (h *Hub) closeAllActiveClients() {
	log.Printf("Closing active clients")
	for client := range h.clients {
		close(client.send)
		delete(h.clients, client)
	}
}
func (h *Hub) Run() {
	log.Printf("Running Hub")
	defer log.Printf("Ending Hub")
	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel
	for {
		select {
		case <-ctx.Done():
			h.closeAllActiveClients()
			return
		case client := <-h.register:
			log.Printf("Register new client")
			h.clients[client] = true
		case client := <-h.unregister:
			log.Printf("UnRegister new client")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			//log.Printf("Message arrived")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					log.Printf("Unknown error, closing Client")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
