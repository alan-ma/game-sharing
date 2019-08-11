package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	// Create game servers
	newGameServer := InitializeNewGameServer(0)

	hub := NewHub(newGameServer)
	go hub.processIO()

	// Simple file server
	fileServer := http.FileServer(http.Dir("../public"))
	http.Handle("/", fileServer)

	// Configure websocket route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(hub, w, r)
	})

	// Start the server using the address specified and log errors
	log.Println("http server stated on", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
