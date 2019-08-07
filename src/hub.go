package main

import (
	"time"
)

// Hub represents a live game being played by one or more players
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Inbound display data from game server
	displayData chan []byte

	// new input
	newInput []byte
}

func mockServerReturn(hub *Hub) {
	for {
		select {
		case hub.displayData <- hub.newInput:
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func newHub() *Hub {
	newHub := &Hub{
		broadcast:   make(chan []byte), // TODO: of type InputData
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		displayData: make(chan []byte),
		newInput:    []byte{'h', 'e', 'l', 'l', 'o', '\n'},
	}
	go mockServerReturn(newHub)
	return newHub
}

func (hub *Hub) runGame() {
	for {
		select {
		case client := <-hub.register:
			// Register the client coming from the channel
			hub.clients[client] = true
		case client := <-hub.unregister:
			// Unregister the client and delete from the active list
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		case inputData := <-hub.broadcast:
			// TODO: send inputData and liveGameState to gameServer, get its displayData and gameState
			hub.newInput = inputData
		case outputData := <-hub.displayData:
			// Process each client
			for client := range hub.clients {
				select {
				case client.send <- outputData:
				default:
					// If nothing can be sent, assume the client is dead or stuck
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
