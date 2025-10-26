package sqlite

import (
	"github.com/thetnaingtn/dirty-hand/store"
)

// Ensure DB implements the store.Driver interface
var _ store.Driver = (*DB)(nil)

// Additional methods that might be needed for the Driver interface
// can be added here as the interface evolves
