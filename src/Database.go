package main

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

// "Table" Descriptions:
// UserStates stores a list of states for each user for a specific game

// NewPool returns a pool of connections to Redis
func NewPool(addr string) *redis.Pool {
	return &redis.Pool{
		// Max number of idle connections in the pool.
		MaxIdle: 3,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func getUserStatesObjectPrefix(gameID GameID, userID UserID) string {
	return "table: UserStates, gameID: " + gameID + ", userID: " + userID
}

// AddToUserStates adds a state model to a user's list of saved states
func AddToUserStates(gameID GameID, userID UserID, newState *State) error {
	conn := DatabasePool.Get()
	defer conn.Close()

	key := getUserStatesObjectPrefix(gameID, userID)

	// Read value from database
	storedValue, readErr := redis.String(conn.Do("GET", key))
	if readErr != nil && readErr != redis.ErrNil {
		// If saved states do not exist, it is empty
		return readErr
	}

	stateList := []State{}

	decodeErr := json.Unmarshal([]byte(storedValue), &stateList)
	if decodeErr != nil {
		return decodeErr
	}

	stateList = append(stateList, *newState)

	// Serialize state list to json
	jsonValue, encodeErr := json.Marshal(stateList)
	if encodeErr != nil {
		return encodeErr
	}

	// Store updated value in database
	_, writeErr := conn.Do("SET", key, jsonValue)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

// GetUserStates returns the saved states that a user has from the database
func GetUserStates(gameID GameID, userID UserID) (string, error) {
	conn := DatabasePool.Get()
	defer conn.Close()

	key := getUserStatesObjectPrefix(gameID, userID)

	// Read value from database
	storedValue, readErr := redis.String(conn.Do("GET", key))

	if readErr != nil && readErr == redis.ErrNil {
		// Not an error if empty
		return "", nil
	}

	return storedValue, readErr
}
