package db

import (
	"github.com/thetnaingtn/dirty-hand/internal/config"
	"github.com/thetnaingtn/dirty-hand/store"
	"github.com/thetnaingtn/dirty-hand/store/db/sqlite"
)

var _ store.Driver = (*sqlite.DB)(nil)

func NewDBDriver(config *config.Config) (store.Driver, error) {
	db, err := sqlite.NewDB(config)
	if err != nil {
		return nil, err
	}

	return db, nil
}
