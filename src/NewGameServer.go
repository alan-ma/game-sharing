package main

var _ GameServer = (*NewGameServer)(nil)
var _ GameState = (*NewGameState)(nil)

// NewGameServer is a concrete instance of GameServer
type NewGameServer struct {
	newestStateID SafeStateID
	savedStates   map[StateID]NewGameState
}

// NewGameState is a concrete instance of GameState
type NewGameState struct {
	id           StateID
	gameServerID GameServerID
	newestText   []byte
	displayData  DisplayData
}

// ProcessState updates the GameState along with new DisplayData based on InputData
func (server *NewGameServer) ProcessState(state GameState, inputs InputData) {
	newState := state.(*NewGameState)
	newState.newestText = append(newState.newestText, inputs...)
	newState.newestText = newState.newestText[:8]
}

// SaveState overwrites the state of the game in the database if it exists, otherwise calls SaveAs
func (server *NewGameServer) SaveState(state GameState) {
	newState := state.(*NewGameState)

	// Check that the current game state exists in the database
	_, ok := server.savedStates[newState.GetID()]
	if !ok {
		panic("Current state does not exist in database")
	}

	// Overwrite existing state in database
	server.savedStates[newState.GetID()] = *newState
}

// SaveAsState saves the state of the game with a new state id
func (server *NewGameServer) SaveAsState(state GameState) {
	newState := state.(*NewGameState)

	// Get a new id and insert it into the database
	newStateID := server.newestStateID.GetAndIncrementSafeStateID()

	// This new id should not exist in the database
	_, ok := server.savedStates[newStateID]
	if ok {
		panic("New state id already exists")
	}

	// Reset state id
	newState.id = newStateID

	// Save in database
	server.savedStates[newStateID] = *newState
}

// LoadState retrieves the GameState from the database
func (server *NewGameServer) LoadState(stateID StateID) GameState {
	newGameState := server.savedStates[stateID]
	return &newGameState
}

// NewState returns a new initialized GameState
func (server *NewGameServer) NewState() GameState {
	newGameState := make(NewGameState)
	return &newGameState
}

// GetID returns the StateID used to access the state
func (state *NewGameState) GetID() StateID {
	return state.id
}

// GetServerID returns the GameServerID, essentially the game the state is for
func (state *NewGameState) GetServerID() GameServerID {
	return state.gameServerID
}

// GetDisplayData returns the current display data for the state
func (state *NewGameState) GetDisplayData() DisplayData {
	return state.displayData
}
