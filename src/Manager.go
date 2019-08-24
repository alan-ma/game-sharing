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
	ID         string `json:"id"`
	LastPlayed string `json:"lastPlayed"`
}

// Games stores a list of the available games to play
var Games []Game

// GameServerMap maps each game ID to its respective server
var GameServerMap map[GameID]GameServer

// UserStates stores a list of states for each user for a specific game
var UserStates map[GameID]map[UserID][]StateID

// Users is a set of existing users
var Users map[UserID]bool

// GetGames returns an index of available games
func GetGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Games)
}

// GetStates returns an index of saved states for a user in a specific game
func GetStates(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]

	if _, ok := GameServerMap[gameID]; !ok {
		http.Error(w, "Game ID does not exist.", http.StatusNotFound)
		return
	}

	userID := params["userID"]
	if _, ok := Users[userID]; !ok {
		http.Error(w, "User ID does not exist.", http.StatusNotFound)
		return
	}

	if _, ok := UserStates[userID]; !ok {
		UserStates[gameID][userID] = make([]StateID, 0)
	}

	states := make([]State, 0)

	for _, stateID := range UserStates[gameID][userID] {
		states = append(states, State{
			ID:         strconv.Itoa(stateID),
			LastPlayed: GameServerMap[gameID].LoadState(stateID).GetSavedDate().String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserStates[gameID][userID])
}

// CreateState starts a new game session
func CreateState(w http.ResponseWriter, r *http.Request) {

}

// LoadState loads a saved state as a live game session for a user
func LoadState(w http.ResponseWriter, r *http.Request) {

}

// SaveState saves the live games session as a saved state
func SaveState(w http.ResponseWriter, r *http.Request) {

}

// Login ensures that the userID exists
func Login(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	if _, ok := Users[userID]; !ok {
		Users[userID] = true
	}
}
