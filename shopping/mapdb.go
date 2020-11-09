package shopping

/**
	Generic Database implementation using an underlying Golang Map
	The Map can be of any shape e.g. map[string]int, map[int]interface{}
	(You can even use a *sqlx.DB instead for a SQL like DB impl, the mutex might not be required)
	which can be determined by the app owner on the data type they want to store.
	This is more of a design choice.

	sync RW mutex has been used to avoid race conditions during concurrent reads and writes.
*/

import (
	"sync"
)

// MapDb Generic DB implementation with mutex
type MapDb struct {
	// mutex for sync RW access
	mu sync.RWMutex
	// generic map instance inside db
	mp interface{}
}

// NewMapDb to create a new MapDb instance
func NewMapDb(data interface{}) *MapDb {
	return &MapDb{
		mp: data,
	}
}