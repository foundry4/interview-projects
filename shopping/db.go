package shopping

/**
Generic Database implementation using a Type which could be anything
The Type can be of any shape e.g. map[string]int, map[int]interface{} as we used Golang map in our app
(You can even use a *sqlx.DB instead for a SQL like DB impl, the mutex might not be required)
which can be determined by the app owner on the data type they want to store.
This is more of a design choice.

sync RW mutex has been used to avoid race conditions during concurrent reads and writes.
*/

import (
	"sync"
)

// DB Generic DB implementation with mutex
type DB struct {
	// mutex for sync RW access
	Mu sync.RWMutex
	// generic map instance inside db
	Type interface{}
}

// NewDB to create a new MapDb instance
func NewDB(data interface{}) *DB {
	return &DB{
		Type: data,
	}
}
