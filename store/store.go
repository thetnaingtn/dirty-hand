package store

import "github.com/thetnaingtn/dirty-hand/internal/config"

type Store struct {
	driver Driver
	config *config.Config
}

func NewStore(driver Driver, config *config.Config) *Store {
	return &Store{
		driver: driver,
		config: config,
	}
}

func (s *Store) Close() error {
	if s.driver == nil {
		return nil
	}

	return s.driver.Close()
}
