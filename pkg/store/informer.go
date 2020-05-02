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
	// RegisteredInformers contains a list of store Informers configured using the package
	// it is useful for printing debug information related to store.
	RegisteredInformers map[string]*Informer = make(map[string]*Informer)

	// ControllerManager is a global controller manager for all the store Informers
	ControllerManager *controller.Manager = controller.NewManager()
)

// registerInformer adds the specified Informer to registered Informer list.
func registerInformer(c *Informer) {
	log.WithFields(log.Fields{
		"package": "store",
	}).Infof("registering store informer: %s for key: %s", c.name, c.Key)

	RegisteredInformers[c.name] = c
}

type addFuncType func(*v1alpha1.KVPairStruct) error

type updateFuncType func(*v1alpha1.KVPairStruct, uint64) error

type deleteFuncType func(string) error

// Informer is a type corresponding to a store Informer.
// A store Informer can be used to run function based on changes to a store
// object.
type Informer struct {
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
func (c *Informer) Name() string {
	return c.name
}

// NewInformer returns a new store controller to periodically run functions
// based on changes to the specifed key in the store.
func NewInformer(
	key string,
	addFunc func(*v1alpha1.KVPairStruct) error,
	delFunc func(string) error,
	updateFunc func(*v1alpha1.KVPairStruct, uint64) error) *Informer {

	name := utils.RandToken(defaults.StoreControllerNameLength)

	return &Informer{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: delFunc,
		Key:        key,
		name:       name,
		cache:      NewCache(name),
		Manager:    ControllerManager,
	}
}

// NewInformerWithSharedCache returns a new store Informer to periodically run functions
// based on changes to the specifed key in the store.
// The Informer in this case is configured with the global shared store, which share its state
// with other Informer.
func NewInformerWithSharedCache(
	key string,
	addFunc func(*v1alpha1.KVPairStruct) error,
	delFunc func(string) error,
	updateFunc func(*v1alpha1.KVPairStruct, uint64) error) *Informer {

	return &Informer{
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
func (c *Informer) Run() error {
	log.WithFields(log.Fields{
		"package": "store",
	}).Infof("starting to setup store informer: %s", c.name)
	registerInformer(c)

	fun, err := controller.NewControllerFunction(c.storeInformerDoFunc)
	if err != nil {
		return fmt.Errorf("error while creating informer's controller function: %s", err)
	}
	// create controller for the perodically running function.
	err = c.Manager.UpdateController(
		c.name,
		"store-informer",
		controller.Internal{
			DoFunc:      fun,
			RunInterval: defaults.StoreControllerRunInterval,
		})
	if err != nil {
		return fmt.Errorf("error while updating controller: %s", err)
	}

	log.Infof("starting to run, store cache controller function for deleted keys.")
	c.cache.RegisterDeleteFunc(c.Key, c.DeleteFunc)
	err = c.cache.RunController()
	if err != nil {
		return err
	}
	return nil
}

func (c *Informer) storeInformerDoFunc(ctx context.Context) error {
	log.WithFields(log.Fields{
		"package": "store",
	}).Debugf("running controller function for store controller: %s", c.name)

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

	log.Debugf("informer controller execution finished")
	return nil
}

// DeleteFromCache deletes the entry from the store cache.
func (c *Informer) DeleteFromCache(key string) {
	c.cache.Remove(key)
}

// Stop shuts down a running store controller.
func (c *Informer) Stop() error {
	log.WithFields(log.Fields{
		"package": "store",
	}).Infof("stopping controller: %s", c.name)
	err := c.Manager.RemoveControllerAndWait(c.name)
	if err != nil {
		return err
	}

	c.cache.StopController()
	delete(RegisteredInformers, c.name)
	return nil
}
