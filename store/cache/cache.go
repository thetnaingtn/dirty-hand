package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	DefaultTTL time.Duration
	MaxItems   int
}

type Cache struct {
	data       sync.Map
	config     Config
	itemCount  atomic.Int64
	stopedChan chan struct{}
	closedChan chan struct{}
}

type item struct {
	value      any
	expiration time.Time
}

func DefaultConfig() Config {
	return Config{
		DefaultTTL: time.Minute * 10,
		MaxItems:   1000,
	}
}

func NewDefault() *Cache {
	return New(DefaultConfig())
}

func New(config Config) *Cache {
	return &Cache{
		config:     config,
		stopedChan: make(chan struct{}),
		closedChan: make(chan struct{}),
	}
}

func (c *Cache) Set(key string, value any) {
	if _, exist := c.data.Load(key); exist {
		c.data.Delete(key)
	} else {
		c.itemCount.Add(1)
	}

	item := item{
		value:      value,
		expiration: time.Now().Add(c.config.DefaultTTL),
	}

	c.data.Store(key, item)
}

func (c *Cache) Get(key string) (any, bool) {
	value, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	itm, ok := value.(item)
	if !ok {
		c.data.Delete(key)
		c.itemCount.Add(-1)
		return nil, false
	}

	if time.Now().After(itm.expiration) {
		c.data.Delete(key)
		c.itemCount.Add(-1)
		return nil, false
	}

	return itm.value, true
}
