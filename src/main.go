package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	// Initialize games, users, and game servers
	Games = append(Games,
		Game{
			ID:          "0",
			ImageNumber: "0",
			Name:        "New Game",
			Description: "A new game to play!",
		})
	GameServerMap = make(map[GameID]GameServer)
	UserStates = make(map[GameID]map[UserID][]StateID)
	Users = make(map[UserID]bool)
	for _, game := range Games {
		switch game.ID {
		case "0":
			GameServerMap["0"] = InitializeNewGameServer(0)
		}
		UserStates[game.ID] = make(map[UserID][]StateID)
	}

	// This needs to be done in rest api
	// hub := NewHub(newGameServer)
	// go hub.processIO()
	// Configure websocket route
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	ServeWebSocket(hub, w, r)
	// })

	// Simple file server
	fileServer := http.FileServer(http.Dir("../public"))
	http.Handle("/", fileServer)

	// Initialize router
	router := mux.NewRouter()

	// Define RESTful endpoints
	router.HandleFunc("/games", GetGames).Methods("GET")
	router.HandleFunc("/games/{id}/{userID}", GetStates).Methods("GET")
	router.HandleFunc("/games/{id}/{userID}", CreateState).Methods("PUT")
	router.HandleFunc("/games/{id}/{userID}/{StateID}", LoadState).Methods("GET")
	router.HandleFunc("/games/{id}/{userID}", SaveState).Methods("POST")
	router.HandleFunc("/login/{id}", Login).Methods("POST")

	// Start the server using the address specified and log errors
	log.Println("http server stated on", *addr)
	err := http.ListenAndServe(*addr, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
