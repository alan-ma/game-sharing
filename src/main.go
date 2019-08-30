package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "HTTP service address")
var wsAddr = flag.String("wsAddr", ":8082", "WebSocket service address")
var redisAddr = flag.String("redisAddr", ":6379", "Redis service address")

// MainRouter handles the RESTful API endpoints
var MainRouter *mux.Router

// WSRouter handles the websocket connections
var WSRouter *mux.Router

// DatabasePool is the pool of connections to Redis
var DatabasePool *redis.Pool

func main() {
	flag.Parse()

	// Initialize Redis database
	DatabasePool = NewPool(*redisAddr)
	log.Println("Redis server stated on", *redisAddr)

	// Test Redis connection
	conn := DatabasePool.Get()
	defer conn.Close()
	databaseErr := Ping(conn)
	if databaseErr != nil {
		fmt.Println(databaseErr)
	}

	// Initialize games, users, and game servers
	Games = append(Games,
		Game{
			ID:          "0",
			ImageNumber: "0",
			Name:        "New Game",
			Description: "A new game to play!",
		})
	GameServerMap = make(map[GameID]GameServer)
	Users = make(map[UserID]bool)
	Hubs = make(map[StateID]*Hub)
	UserClients = make(map[UserID]*Client)
	for _, game := range Games {
		switch game.ID {
		case "0":
			GameServerMap["0"] = InitializeNewGameServer(0)
		}
	}

	// Initialize router
	MainRouter = mux.NewRouter()
	WSRouter = mux.NewRouter()

	// Simple file server
	fileServer := http.FileServer(http.Dir("../public"))
	MainRouter.Handle("/", fileServer)

	// Define RESTful endpoints
	MainRouter.HandleFunc("/games", GetGames).Methods("GET")
	MainRouter.HandleFunc("/games/{id}/{userID}", GetStates).Methods("GET")
	MainRouter.HandleFunc("/games/{id}/{userID}", CreateState).Methods("PUT")
	MainRouter.HandleFunc("/games/{id}/{userID}/{stateID}", LoadState).Methods("GET")
	MainRouter.HandleFunc("/games/{id}/{userID}/{stateID}", SaveState).Methods("PUT")
	MainRouter.HandleFunc("/login/{id}", Login).Methods("POST")

	// Configure websocket route
	WSRouter.HandleFunc("/play/{id}/{userID}/{stateID}", HandleWebSocket)

	// Start the server using the address specified and log errors
	log.Println("HTTP server stated on", *addr)
	go func() {
		err := http.ListenAndServe(*addr, MainRouter)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	// Start the websocket connection using the address specified and log errors
	log.Println("WebSocket server stated on", *wsAddr)
	err := http.ListenAndServe(*wsAddr, WSRouter)
	if err != nil {
		log.Fatal("ListenAndServe (Websocket): ", err)
	}
}
