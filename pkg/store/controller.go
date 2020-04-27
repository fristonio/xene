package store

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/fristonio/xene/pkg/controller"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/fristonio/xene/pkg/utils"
)

var (
	// RegisteredControllers contains a list of store controllers configured using the package
	// it is useful for printing debug information related to store.
	RegisteredControllers map[string]*Controller = make(map[string]*Controller)

	// ControllerManager is a global controller manager for all the store controllers
	ControllerManager *controller.Manager = controller.NewManager()
)

// registerController adds the specified controller to registered controller list.
func registerController(c *Controller) {
	log.WithFields(log.Fields{
		"package": "store",
	}).Infof("registering store controller: %s for key: %s", c.name, c.Key)

	RegisteredControllers[c.name] = c
}

type addFuncType func(*v1alpha1.KVPairStruct) error

type updateFuncType func(*v1alpha1.KVPairStruct, uint64) error

type deleteFuncType func(string) error

// Controller is a type corresponding to a store controller.
// A store controller can be used to run function based on changes to a store
// object.
type Controller struct {
	AddFunc addFuncType

	UpdateFunc updateFuncType

	DeleteFunc deleteFuncType

	// Key is the prefix key to watch for.
	Key string

	name string

	cache Cache

	Manager *controller.Manager
}

// Name returns the name of the controller.
func (c *Controller) Name() string {
	return c.name
}

// NewController returns a new store controller to periodically run functions
// based on changes to the specifed key in the store.
func NewController(
	key string,
	addFunc func(*v1alpha1.KVPairStruct) error,
	delFunc func(string) error,
	updateFunc func(*v1alpha1.KVPairStruct, uint64) error) *Controller {

	name := utils.RandToken(defaults.StoreControllerNameLength)

	return &Controller{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: delFunc,
		Key:        key,
		name:       name,
		cache:      NewCache(name),
		Manager:    ControllerManager,
	}
}

// NewControllerWithSharedCache returns a new store controller to periodically run functions
// based on changes to the specifed key in the store.
// The controller in this case is configured with the global shared store, which share its state
// with other controller.
func NewControllerWithSharedCache(
	key string,
	addFunc func(*v1alpha1.KVPairStruct) error,
	delFunc func(string) error,
	updateFunc func(*v1alpha1.KVPairStruct, uint64) error) *Controller {

	return &Controller{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: delFunc,
		Key:        key,
		name:       utils.RandToken(defaults.StoreControllerNameLength),
		cache:      GlobalSharedCache,
		Manager:    ControllerManager,
	}
}

// Run starts running the store controller configured.
func (c *Controller) Run() error {
	log.WithFields(log.Fields{
		"package": "store",
	}).Infof("starting to setup store controller: %s", c.name)
	registerController(c)

	fun, err := controller.NewControllerFunction(c.storeControllerDoFunc)
	if err != nil {
		return fmt.Errorf("error while creating controller function: %s", err)
	}
	// create controller for the perodically running function.
	err = c.Manager.UpdateController(
		c.name,
		"Store",
		controller.Internal{
			DoFunc:      fun,
			RunInterval: defaults.StoreControllerRunInterval,
		})
	if err != nil {
		return fmt.Errorf("error while updating controller: %s", err)
	}

	log.Infof("starting to run, store cache controller function for deleted keys.")
	c.cache.RegisterDeleteFunc(c.Key, c.DeleteFunc)
	_ = c.cache.RunController()
	return nil
}

func (c *Controller) storeControllerDoFunc(ctx context.Context) error {
	log.WithFields(log.Fields{
		"package": "store",
	}).Infof("running controller function for store controller: %s", c.name)

	// Here we do a prefix scan for the provided key with our function
	KVStore.PrefixScanWithFunction(context.TODO(), c.Key, func(kv *v1alpha1.KVPairStruct) {
		if c.cache.CheckIfExists(kv.Key) && kv.DeletedOrExpired {
			err := c.DeleteFunc(kv.Key)
			if err != nil {
				log.Errorf("error while calling delete for key: %s: %s", kv.Key, err)
				return
			}
			c.cache.Remove(kv.Key)
			return
		}

		if c.cache.CheckIfAdded(kv.Key, kv.Version) {
			err := c.AddFunc(kv)
			if err != nil {
				log.Errorf("error while calling add for key: %s: %s", kv.Key, err)
				return
			}
		} else if c.cache.CheckIfUpdated(kv.Key, kv.Version) {
			err := c.UpdateFunc(kv, c.cache.Get(kv.Key))
			if err != nil {
				log.Errorf("error while calling update for key: %s: %s", kv.Key, err)
				return
			}
		}

		c.cache.Set(kv.Key, kv.Version)
	})

	log.Infof("controller execution finished")
	return nil
}

// DeleteFromCache deletes the entry from the store cache.
func (c *Controller) DeleteFromCache(key string) {
	c.cache.Remove(key)
}

// Stop shuts down a running store controller.
func (c *Controller) Stop() error {
	log.WithFields(log.Fields{
		"package": "store",
	}).Infof("stopping controller: %s", c.name)
	err := c.Manager.RemoveControllerAndWait(c.name)
	if err != nil {
		return err
	}

	c.cache.StopController()
	delete(RegisteredControllers, c.name)
	return nil
}
