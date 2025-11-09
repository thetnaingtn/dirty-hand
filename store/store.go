package store

import (
	"github.com/thetnaingtn/dirty-hand/internal/config"
	"github.com/thetnaingtn/dirty-hand/store/cache"
)

type Store struct {
	driver       Driver
	config       *config.Config
	sessionCache *cache.Cache
	userCache    *cache.Cache
}

func NewStore(driver Driver, config *config.Config) *Store {
	return &Store{
		driver:       driver,
		config:       config,
		sessionCache: cache.NewDefault(),
		userCache:    cache.NewDefault(),
	}
}

func (s *Store) Close() error {
	if s.driver == nil {
		return nil
	}

	return s.driver.Close()
}
