package product

import "sync"

func NewCache() *cache {
	return &cache{
		types: make(map[string]int, 3),
		mux:   &sync.RWMutex{},
	}
}

type cache struct {
	types map[string]int
	mux   *sync.RWMutex
}

func (c *cache) Get(name string) (int, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	id, ok := c.types[name]
	return id, ok
}

func (c *cache) Set(name string, id int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.types[name] = id
}
