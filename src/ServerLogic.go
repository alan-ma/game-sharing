package main

// ServerLogic provides the shared functionality between game servers
type ServerLogic struct {
	newestStateID SafeStateID
	savedStates   map[StateID]GameState
}

// SaveState overwrites the state of the game in the database if it exists, otherwise calls SaveAs
func (server *ServerLogic) SaveState(state GameState) {
	// Check that the current game state exists in the database
	_, ok := server.savedStates[state.GetID()]
	if !ok {
		panic("Current state does not exist in database")
	}

	// Overwrite existing state in database
	server.savedStates[state.GetID()] = state
}

// SaveAsState saves the state of the game with a new state id
func (server *ServerLogic) SaveAsState(state GameState) {
	// Get a new id and insert it into the database
	newStateID := server.newestStateID.GetAndIncrementSafeStateID()

	// This new id should not exist in the database
	_, ok := server.savedStates[newStateID]
	if ok {
		panic("New state id already exists")
	}

	// Reset state id
	state.SetID(newStateID)

	// Save in database
	server.savedStates[newStateID] = state
}

// LoadState retrieves the GameState from the database
func (server *ServerLogic) LoadState(stateID StateID) GameState {
	newGameState := server.savedStates[stateID]
	return newGameState
}
