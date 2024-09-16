package gonvoy

import (
	"sync"
)

// Cache is an interface that defines methods for storing and retrieving data in an internal cache.
// It is designed to maintain data persistently throughout Envoy's lifespan.
type Cache interface {
	// Store allows you to save a value of any type under a key of any type.
	//
	// Please use caution! The Store function overwrites any existing data.
	Store(key, value any)
}

type inmemoryCache struct {
	stash sync.Map
}

func newInternalCache() *inmemoryCache {
	return &inmemoryCache{}
}

func (c *inmemoryCache) Store(key, value any) {
	c.stash.Store(key, value)
}

// LoadValue retrieves a value associated with a specific key from the cache
// and returns it as the specified type T.
//
// It returns the value of type T, a boolean indicating whether the value
// was found, and an error if an incompatible type is received.
//
// The function performs a type-safe cast to the specified type T, and
// if the stored value cannot be cast to T, it returns ErrIncompatibleReceiver.
//
// Example usage:
//
//	type mystruct struct{}
//
//	data := mystruct{}
//	cache.Store("keyName", data)
//
//	result, found, err := LoadValue[mystruct](cache, "keyName")
//	if err != nil {
//	    // Handle error (e.g., ErrIncompatibleReceiver)
//	}
//	if found {
//	    // Use result (which is of type mystruct)
//	}
func LoadValue[T any](c *inmemoryCache, key interface{}) (T, bool, error) {
	var zero T

	v, ok := c.stash.Load(key)
	if !ok {
		return zero, false, nil
	}

	src, ok := v.(T)
	if !ok {
		return zero, false, ErrIncompatibleReceiver
	}

	return src, true, nil
}
