package store

// Cache is a map containing store keys corresponding to the last revision applied for the
// same object.
type Cache map[string]uint64

// NewCache returns a new instance of Store cache
func NewCache() Cache {
	return make(Cache)
}

// CheckIfUpdated checks if the key and the revision provided have been added or
// updated looking at the configured cache.
func (c Cache) CheckIfUpdated(key string, revision uint64) bool {
	if lastRevision, ok := c[key]; ok {
		if lastRevision < revision {
			return true
		}
	}

	return false
}

// CheckIfAdded checks if the key and the revision provided have been added
// looking at the configured cache.
func (c Cache) CheckIfAdded(key string, revision uint64) bool {
	if _, ok := c[key]; !ok {
		return true
	}

	return false
}

// CheckIfExists checks if the provided key exists in the cache or not.
func (c Cache) CheckIfExists(key string) bool {
	if _, ok := c[key]; ok {
		return true
	}

	return false
}

// Set adds or update an existing key in the Cache
func (c Cache) Set(key string, revision uint64) {
	c[key] = revision
}

// Get returns the saved version of the provided key in the cache.
func (c Cache) Get(key string) uint64 {
	return c[key]
}

// Remove deletes the key from the Cache.
func (c Cache) Remove(key string) {
	if _, ok := c[key]; ok {
		delete(c, key)
	}
}

// GlobalSharedCache is the global shared cache that can be used
// by any of the store controller.
var GlobalSharedCache Cache = make(Cache)
