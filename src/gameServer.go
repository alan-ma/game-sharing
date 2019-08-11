package main

import (
	"sync"
)

// GameServer is an interface for the main actions of the game
type GameServer interface {
	// Run the game logic given certain inputs
	// The GameState will be updated along with new DisplayData
	ProcessState(GameState, InputData)

	// Save and load the game state
	SaveState(GameState)
	SaveAsState(GameState)
	LoadState(StateID) GameState
	NewState() GameState
}

// GameState holds the information needed by the game
type GameState interface {
	GetID() StateID
	SetID(StateID)
	GetServerID() GameServerID
	GetDisplayData() DisplayData
}

// StateID is a generated unique id for each GameState
type StateID int

func (id *StateID) increment() {
	*id++
}

// SafeStateID is a StateID safe to use concurrently, prevents race condition
type SafeStateID struct {
	id  StateID
	mux sync.Mutex
}

// GameServerID is a constant that identifies which game is being played
type GameServerID int

// GetAndIncrementSafeStateID returns the current value of the SafeStateID and increments it after
func (safeStateID *SafeStateID) GetAndIncrementSafeStateID() StateID {
	// Lock so only one goroutine at a time can access the StateID
	safeStateID.mux.Lock()
	newStateID := safeStateID.id
	safeStateID.id.increment()
	safeStateID.mux.Unlock()
	return newStateID
}
