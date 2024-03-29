package main

import (
	"time"
)

// Hub represents a live game being played by one or more players
type Hub struct {
	// Game server that processes the game state
	server GameServer

	// The current game state
	state GameState

	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan InputData

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Inbound display data from game server
	displayData chan DisplayData

	// Game input data
	gameInput InputData
}

func runGameLoop(hub *Hub) {
	for {
		hub.server.ProcessState(hub.state, hub.gameInput)
		hub.displayData <- hub.state.GetDisplayData()
		hub.gameInput = nil
		time.Sleep(10 * time.Millisecond) // probably some other way to make a consistent loop
	}
}

// NewHub returns a new Hub for the live game
func NewHub(server GameServer) *Hub {
	newHub := &Hub{
		server:      server,
		state:       server.NewState(),
		broadcast:   make(chan InputData),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		displayData: make(chan DisplayData),
	}
	go runGameLoop(newHub)
	return newHub
}

// LoadHub returns a new Hub with a given state
func LoadHub(server GameServer, stateID StateID) *Hub {
	loadedState := server.LoadState(stateID)
	newHub := &Hub{
		server:      server,
		state:       loadedState,
		broadcast:   make(chan InputData),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		displayData: make(chan DisplayData),
	}
	go runGameLoop(newHub)
	return newHub
}

func (hub *Hub) processIO() {
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
		case newInput := <-hub.broadcast:
			for _, char := range newInput {
				hub.gameInput = append(hub.gameInput, char)
			}
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
