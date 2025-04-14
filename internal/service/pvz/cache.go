package pvz

import "sync"

func NewCache() *cache {
	return &cache{
		cities: make(map[string]int, 3),
		mux:    &sync.RWMutex{},
	}
}

type cache struct {
	cities map[string]int
	mux    *sync.RWMutex
}

func (c *cache) Get(name string) (int, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	id, ok := c.cities[name]
	return id, ok
}

func (c *cache) Set(name string, id int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cities[name] = id
}
