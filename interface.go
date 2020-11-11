package cus

import "errors"

var (
	ErrDuplicateObject = errors.New("Duplicate object")
	ErrNotFound        = errors.New("Not found")
)

// Unique store represents the interface for storing and retrieving unique
// objects from/to database
type UniqueStore interface {
	// Create unique object and store it in cassandra
	// In case of the object data is already persisted with other object
	// the function returns the `duplicateId` of the object along with err = `ErrDuplicateObject`
	Create(id string, data []byte) (duplicateId *string, err error)
	// Get returns the object data
	// In case the object is not found - then err = `ErrNotFound` and the data will be `nil`
	Get(id string) (data []byte, err error)
	// CreateSchema creates the tables if they do not exist in the default cluster keyspace
	// If `force` flag is set to true, then the table will be first droped if it exists
	CreateSchema(force bool) error
	// Delete removes the data from the db
	// In case the object is not found - then err = `ErrNotFound` and the data will be `nil`
	// In any other case the original error will be propagated
	Delete(id string) (err error)
}
