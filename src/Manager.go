package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GameID is the identifier for each game
type GameID = string

// UserID is the identifier for each user
type UserID = string

// Game is the model for game information
type Game struct {
	ID          GameID `json:"id"`
	ImageNumber string `json:"imageNumber"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// State is the model for state information
type State struct {
	ID      string `json:"id"`
	SavedOn string `json:"savedOn"`
}

// Games stores a list of the available games to play
var Games []Game

// GameServerMap maps each game ID to its respective server
var GameServerMap map[GameID]GameServer

// UserStates stores a list of states for each user for a specific game
var UserStates map[GameID]map[UserID][]StateID

// Users is a set of existing users
var Users map[UserID]bool

// Hubs is a map of live game sessions
var Hubs map[StateID]*Hub

// UserClients is a map of users to their clients
var UserClients map[UserID]*Client

func errorCheck(w http.ResponseWriter, r *http.Request, gameID GameID, userID UserID) bool {
	if _, ok := GameServerMap[gameID]; !ok {
		http.Error(w, "Game ID does not exist.", http.StatusNotFound)
		return false
	}

	if _, ok := Users[userID]; !ok {
		http.Error(w, "User ID does not exist.", http.StatusNotFound)
		return false
	}

	if _, ok := UserStates[userID]; !ok {
		UserStates[gameID][userID] = make([]StateID, 0)
	}

	return true
}

// getValidStateID returns a state ID if valid, otherwise -1
func getValidStateID(w http.ResponseWriter, r *http.Request, stateIDStr string) StateID {
	stateID, err := strconv.Atoi(stateIDStr)
	if err != nil {
		http.Error(w, "Invalid game session ID.", http.StatusBadRequest)
		return -1
	}
	return stateID
}

// isValidLiveSession checks if the state ID given is a valid live game session
func isValidLiveSession(w http.ResponseWriter, r *http.Request, stateID StateID) bool {
	if _, ok := Hubs[stateID]; !ok {
		http.Error(w, "State ID is not a valid live game session.", http.StatusNotFound)
		return false
	}

	return true
}

// IsValidSavedState checks if the state ID is saved in the database
func isValidSavedState(w http.ResponseWriter, r *http.Request, stateID StateID) bool {
	return true
}

// GetGames returns an index of available games
func GetGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Games)
}

// GetStates returns an index of saved states for a user in a specific game
func GetStates(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]
	userID := params["userID"]

	if !errorCheck(w, r, gameID, userID) {
		return
	}

	states := make([]State, 0)

	for _, stateID := range UserStates[gameID][userID] {
		states = append(states, State{
			ID:      strconv.Itoa(stateID),
			SavedOn: GameServerMap[gameID].LoadState(stateID).GetSavedDate().String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserStates[gameID][userID])
}

// CreateState starts a new game session and returns the new state ID
func CreateState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]
	userID := params["userID"]

	if !errorCheck(w, r, gameID, userID) {
		return
	}

	// Create a client and hub to handle the websocket connection
	hub := NewHub(GameServerMap[gameID])
	Hubs[hub.state.GetID()] = hub

	// Return the state information to the client
	newState := State{
		ID:      strconv.Itoa(hub.state.GetID()),
		SavedOn: hub.state.GetSavedDate().String(),
	}

	// Start processing I/O on the game hub
	go hub.processIO()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newState)
}

// LoadState loads a saved state as a live game session for a user
func LoadState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]
	userID := params["userID"]
	stateIDStr := params["stateID"]

	if !errorCheck(w, r, gameID, userID) {
		return
	}

	stateID := getValidStateID(w, r, stateIDStr)
	if stateID == -1 || !isValidLiveSession(w, r, stateID) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(GameServerMap[gameID].GetServerLogic().savedStates[stateID]))
}

// SaveState saves the live games session as a saved state
func SaveState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]
	userID := params["userID"]
	stateIDStr := params["stateID"]

	if !errorCheck(w, r, gameID, userID) {
		return
	}

	stateID := getValidStateID(w, r, stateIDStr)
	if stateID == -1 || !isValidLiveSession(w, r, stateID) {
		return
	}

	// Save the state to the database
	newStateID, savedOn := GameServerMap[gameID].SaveAsState(stateID)

	// Return the state information to the client
	newState := State{
		ID:      strconv.Itoa(newStateID),
		SavedOn: savedOn.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newState)
}

// Login ensures that the userID exists
func Login(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	if _, ok := Users[userID]; !ok {
		Users[userID] = true
		w.WriteHeader(http.StatusCreated)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// HandleWebSocket processes a request for a WebSocket
// i.e. connects a client to a game hub running a live game session
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]
	userID := params["userID"]
	stateIDStr := params["stateID"]

	if !errorCheck(w, r, gameID, userID) {
		return
	}

	// TODO: Change to actual auth
	if _, ok := Users[userID]; !ok {
		http.Error(w, "Not authorized.", http.StatusUnauthorized)
		return
	}

	stateID := getValidStateID(w, r, stateIDStr)
	if stateID == -1 || !isValidLiveSession(w, r, stateID) {
		return
	}

	ServeWebSocket(userID, stateID, w, r)
}
