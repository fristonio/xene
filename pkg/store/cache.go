package store

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fristonio/xene/pkg/controller"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Cache is a map containing store keys corresponding to the last revision applied for the
// same object.
type Cache struct {
	cache   map[string]uint64
	manager *controller.Manager
	name    string

	deleteFuncs map[string]deleteFuncType
	mux         *sync.Mutex
}

// NewCache returns a new instance of Store cache
func NewCache(name string) Cache {
	return Cache{
		cache:       make(map[string]uint64),
		manager:     controller.NewManager(),
		name:        name,
		deleteFuncs: make(map[string]deleteFuncType),
		mux:         &sync.Mutex{},
	}
}

// CheckIfUpdated checks if the key and the revision provided have been added or
// updated looking at the configured cache.
func (c *Cache) CheckIfUpdated(key string, revision uint64) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if lastRevision, ok := c.cache[key]; ok {
		if lastRevision < revision {
			return true
		}
	}

	return false
}

// CheckIfAdded checks if the key and the revision provided have been added
// looking at the configured cache.
func (c *Cache) CheckIfAdded(key string, revision uint64) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.cache[key]; !ok {
		return true
	}

	return false
}

// CheckIfExists checks if the provided key exists in the cache or not.
func (c *Cache) CheckIfExists(key string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.cache[key]; ok {
		return true
	}

	return false
}

// Set adds or update an existing key in the Cache
func (c *Cache) Set(key string, revision uint64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache[key] = revision
}

// Get returns the saved version of the provided key in the cache.
func (c *Cache) Get(key string) uint64 {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.cache[key]
}

// Remove deletes the key from the Cache.
func (c *Cache) Remove(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.cache, key)
}

// StopController stops the running controller.
func (c *Cache) StopController() {
	err := c.manager.RemoveControllerAndWait(c.name)
	if err != nil {
		log.Errorf("error while stopping cache controller(%s): %s", c.name, err)
	}
}

// RegisterDeleteFunc registers a delete function for delete on the registered
// prefix.
func (c *Cache) RegisterDeleteFunc(prefix string, deleteFunc deleteFuncType) {
	c.deleteFuncs[prefix] = deleteFunc
}

// RunController starts running the cache controller.
func (c *Cache) RunController() error {
	if c.manager.Lookup(c.name) != nil {
		log.Infof("cache controller already running.")
		return nil
	}
	fun, err := controller.NewControllerFunction(c.cacheControllerDoFunc)
	if err != nil {
		return fmt.Errorf("error while creating controller function: %s", err)
	}
	// create controller for the perodically running function.
	err = c.manager.UpdateController(
		c.name,
		"cache",
		controller.Internal{
			DoFunc:      fun,
			RunInterval: defaults.StoreControllerRunInterval,
		})
	if err != nil {
		return fmt.Errorf("error while updating controller: %s", err)
	}
	log.Infof("started running store cache controller")

	return nil
}

func (c *Cache) cacheControllerDoFunc(_ctx context.Context) error {
	log.Infof("running cache controller do func.")
	mErr := errors.NewMultiError()
	for key := range c.cache {
		errs := errors.NewMultiError()
		ok, err := KVStore.Exists(context.TODO(), key)
		if err != nil {
			errs.Append(fmt.Errorf("error while fetching key(%s) from kvstore: %s", key, err))
			continue
		}

		if ok {
			continue
		}

		log.Debugf("key %s has been deleted.", key)
		for prefix, fun := range c.deleteFuncs {
			if strings.HasPrefix(key, prefix) {
				err = fun(key)
				if err != nil {
					errs.Append(fmt.Errorf("error while executing delete func: %s", err))
					continue
				}
			}
		}

		if !errs.HasErrors() {
			c.mux.Lock()
			delete(c.cache, key)
			c.mux.Unlock()
			mErr.Append(errs.GetError())
		}
	}

	return mErr.GetError()
}

// GlobalSharedCache is the global shared cache that can be used
// by any of the store controller.
var GlobalSharedCache Cache = NewCache("global-cache")
