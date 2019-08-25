package main

import (
	"encoding/json"
	"time"
)

// ServerLogic provides the shared functionality between game servers
type ServerLogic struct {
	gameID        GameID
	newestStateID SafeStateID
	savedStates   map[StateID]SavedState
}

// SavedState is the model for saved game information (this is arbitrary)
type SavedState = string

// SaveAsState saves a live game session into the database as a new state
func (server *ServerLogic) SaveAsState(stateID StateID) (StateID, time.Time) {
	hub, isLiveGameSession := Hubs[stateID]

	// Check that the game state is a live game session
	// This check should have been completed already
	if !isLiveGameSession {
		panic("State ID is not a live game session")
	}
	state := hub.state

	// Get a new id and insert it into the database
	newStateID := server.newestStateID.GetAndIncrementSafeStateID()

	// This new id should not exist in the database
	_, alreadyInDatabase := server.savedStates[newStateID]
	if alreadyInDatabase {
		panic("New state ID already exists")
	}

	// Save state in database - converts information to json
	currentTime := time.Now()
	state.SetSavedDate(currentTime)
	stateModel, err := json.Marshal(state)

	// There was something wrong with converting the state to json
	if err != nil {
		panic("Error saving state")
	}

	// Saves the json as a string
	server.savedStates[state.GetID()] = string(stateModel)

	// Reset saved date for live session
	state.ResetSavedDate()

	return newStateID, currentTime
}

// LoadState retrieves the GameState from the database
func (server *ServerLogic) LoadState(stateID StateID) GameState {
	savedState := server.savedStates[stateID]

	// Decode json into custom data type
	// Switch to determine which game state to decode into
	var loadedState GameState
	switch server.gameID {
	case "0":
		loadedState = &NewGameState{}
	}
	err := loadedState.UnmarshalJSON([]byte(savedState))
	if err == nil {
		panic("State did not decode correctly")
	}

	return loadedState
}
