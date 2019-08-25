package main

import (
	"encoding/json"
	"time"
)

// These are to check that the implementation of interfaces is correct
// var _ GameServer = (*NewGameServer)(nil)
// var _ GameState = (*NewGameState)(nil)

// NewGameServer is a concrete instance of GameServer
type NewGameServer struct {
	id          GameServerID
	serverLogic ServerLogic
}

// TODO: remove this
func (server *NewGameServer) GetServerLogic() ServerLogic {
	return server.serverLogic
}

// InitializeNewGameServer starts the game server
func InitializeNewGameServer(id GameServerID) GameServer {
	newServer := &NewGameServer{
		id: id,
		serverLogic: ServerLogic{
			"0", // This is the hard-coded game ID
			SafeStateID{id: 0},
			make(map[StateID]SavedState),
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
	savedDate      time.Time
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

// SaveAsState saves the state of the game with a new state id
func (server *NewGameServer) SaveAsState(stateID StateID) (StateID, time.Time) {
	return server.serverLogic.SaveAsState(stateID)
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
		savedDate:      time.Time{},
	}
	return newGameState
}

// NewStateID returns the latest state ID
func (server *NewGameServer) NewStateID() StateID {
	newStateID := server.serverLogic.newestStateID.GetAndIncrementSafeStateID()
	return newStateID
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

// GetSavedDate returns the time the game state was saved
func (state *NewGameState) GetSavedDate() time.Time {
	return state.savedDate
}

// SetSavedDate sets the time the game state was saved
func (state *NewGameState) SetSavedDate(time time.Time) {
	state.savedDate = time
}

// ResetSavedDate sets the time the game state was saved to zero
func (state *NewGameState) ResetSavedDate() {
	state.savedDate = time.Time{}
}

// IsLiveSession returns true if the state has not been saved yet (i.e. it is a live session being played)
func (state *NewGameState) IsLiveSession() bool {
	return state.savedDate.IsZero()
}

// MarshalJSON returns a json encoded version of NewGameState
func (state *NewGameState) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID             int       `json:"id"`
		ServerID       int       `json:"serverID"`
		SpritePosition int       `json:"spritePosition"`
		DisplayData    []byte    `json:"displayData"`
		SavedDate      time.Time `json:"savedDate"`
	}{
		ID:             state.id,
		ServerID:       state.serverID,
		SpritePosition: state.spritePosition,
		DisplayData:    state.displayData,
		SavedDate:      state.savedDate,
	})
}

// UnmarshalJSON returns a decoded NewGameState
func (state *NewGameState) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID             int       `json:"id"`
		ServerID       int       `json:"serverID"`
		SpritePosition int       `json:"spritePosition"`
		DisplayData    []byte    `json:"displayData"`
		SavedDate      time.Time `json:"savedDate"`
	}{
		ID:             state.id,
		ServerID:       state.serverID,
		SpritePosition: state.spritePosition,
		DisplayData:    state.displayData,
		SavedDate:      state.savedDate,
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	state.id = aux.ID
	state.serverID = aux.ServerID
	state.spritePosition = aux.SpritePosition
	state.displayData = aux.DisplayData
	state.savedDate = aux.SavedDate

	return nil
}
