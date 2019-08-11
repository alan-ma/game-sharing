package main

// These are to check that the implementation of interfaces is correct
// var _ GameServer = (*NewGameServer)(nil)
// var _ GameState = (*NewGameState)(nil)

// NewGameServer is a concrete instance of GameServer
type NewGameServer struct {
	id          GameServerID
	serverLogic ServerLogic
}

// InitializeNewGameServer starts the game server
func InitializeNewGameServer(id GameServerID) GameServer {
	newServer := &NewGameServer{
		id: id,
		serverLogic: ServerLogic{
			SafeStateID{id: 0},
			make(map[StateID]GameState),
		},
	}
	return newServer
}

// NewGameState is a concrete instance of GameState
type NewGameState struct {
	id             StateID
	serverID       GameServerID
	spritePosition int
	displayData    DisplayData
}

// ProcessState updates the GameState along with new DisplayData based on InputData
func (server *NewGameServer) ProcessState(state GameState, inputs InputData) {
	newState := state.(*NewGameState)

	for _, char := range inputs {
		// vbKeyLeft   37  LEFT ARROW key
		// vbKeyUp     38  UP ARROW key
		// vbKeyRight  39  RIGHT ARROW key
		// vbKeyDown   40  DOWN ARROW key
		if char == 37 && newState.spritePosition > 0 {
			// LEFT ARROW key
			newState.spritePosition--
		} else if char == 39 && newState.spritePosition < 7 {
			newState.spritePosition++
		}
	}

	newState.displayData = []byte{48, 48, 48, 48, 48, 48, 48, 48}
	newState.displayData[newState.spritePosition] = 49
}

// SaveState overwrites the state of the game in the database if it exists, otherwise calls SaveAs
func (server *NewGameServer) SaveState(state GameState) {
	server.serverLogic.SaveState(state)
}

// SaveAsState saves the state of the game with a new state id
func (server *NewGameServer) SaveAsState(state GameState) {
	server.serverLogic.SaveAsState(state)
}

// LoadState retrieves the GameState from the database
func (server *NewGameServer) LoadState(stateID StateID) GameState {
	return server.serverLogic.LoadState(stateID)
}

// NewState returns a new initialized GameState
func (server *NewGameServer) NewState() GameState {
	newStateID := server.serverLogic.newestStateID.GetAndIncrementSafeStateID()
	newGameState := &NewGameState{
		id:             newStateID,
		serverID:       server.id,
		spritePosition: 3,
		displayData:    make([]byte, 8),
	}
	return newGameState
}

// GetID returns the StateID used to access the state
func (state *NewGameState) GetID() StateID {
	return state.id
}

// SetID sets the StateID to the new id
func (state *NewGameState) SetID(id StateID) {
	state.id = id
}

// GetServerID returns the GameServerID, essentially the game the state is for
func (state *NewGameState) GetServerID() GameServerID {
	return state.serverID
}

// GetDisplayData returns the current display data for the state
func (state *NewGameState) GetDisplayData() DisplayData {
	return state.displayData
}
