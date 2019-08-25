package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")
var wsAddr = flag.String("wsAddr", ":8082", "websocket service address")

// MainRouter handles the RESTful API endpoints
var MainRouter *mux.Router

// WSRouter handles the websocket connections
var WSRouter *mux.Router

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
	Hubs = make(map[StateID]*Hub)
	for _, game := range Games {
		switch game.ID {
		case "0":
			GameServerMap["0"] = InitializeNewGameServer(0)
		}
		UserStates[game.ID] = make(map[UserID][]StateID)
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
	MainRouter.HandleFunc("/games/{id}/{userID}/{StateID}", LoadState).Methods("GET")
	MainRouter.HandleFunc("/games/{id}/{userID}", SaveState).Methods("POST")
	MainRouter.HandleFunc("/login/{id}", Login).Methods("POST")

	// Configure websocket route
	WSRouter.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		stateID, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "State ID is not a valid live game session.", http.StatusNotFound)
		}
		ServeWebSocket(stateID, w, r)
	})

	// Start the server using the address specified and log errors
	log.Println("http server stated on", *addr)
	go func() {
		err := http.ListenAndServe(*addr, MainRouter)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	// Start the websocket connection using the address specified and log errors
	log.Println("websocket server stated on", *wsAddr)
	err := http.ListenAndServe(*wsAddr, WSRouter)
	if err != nil {
		log.Fatal("ListenAndServe (Websocket): ", err)
	}
}
