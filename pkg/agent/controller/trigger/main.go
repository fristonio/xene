package trigger

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fristonio/xene/pkg/controller"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	log "github.com/sirupsen/logrus"
)

var (
	eventRemove string = "remove"

	eventUpdate string = "update"
)

type storeEvent struct {
	event string
	kv    *v1alpha1.KVPairStruct
}

// Controller contains the implementation of the trigger controller.
type Controller struct {
	// Manager is the controller manager for all the running controller
	// against agent manager.
	Manager *controller.Manager

	// storeInformer is the store informer corresponding to the agent objects
	// in the datastore.
	storeInformer *store.Informer

	// name is the name of the controller, this is derived from the store controller.
	name string

	// Triggers contains the specification of the triggers configured on the agent.
	Triggers map[string]*Trigger

	// DefaultTriggers contains a list of default triggers configured for the agent
	// These are basically Webhook configured triggers and hence don't need a controller
	// to function
	DefaultTriggers map[string]struct{}

	// events contains the list of latest events on Trigger list
	events map[string]storeEvent

	mu *sync.Mutex
}

// TriggerCtrl is the global instance of the trigger controller running for the agent.
var TriggerCtrl *Controller = NewController()

// NewController returns a new agent controller manager which manages all the
// triggers configured on the agent.
func NewController() *Controller {
	return &Controller{
		Manager:         controller.NewManager(),
		Triggers:        make(map[string]*Trigger),
		DefaultTriggers: make(map[string]struct{}),
		events:          make(map[string]storeEvent),
		mu:              &sync.Mutex{},
	}
}

// Type returns the type of controller configured, for trigger this is trigger.
func (c *Controller) Type() string {
	return "Trigger"
}

// Configure sets up the trigger controller and all its required components.
func (c *Controller) Configure() {
	c.storeInformer = c.newTriggerStoreInformer()
	c.name = c.storeInformer.Name()
}

// Run starts running the trigger controller on the agent.
func (c *Controller) Run() error {
	f, err := controller.NewControllerFunction(c.triggerUpdateControllerFunc)
	if err != nil {
		return fmt.Errorf("error while creating update controller: %s", err)
	}

	err = c.Manager.UpdateController("trigger-update", "trigger", controller.Internal{
		DoFunc:      f,
		RunInterval: 15 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error while setting up update controller: %s", err)
	}

	return c.storeInformer.Run()
}

// Stop shuts down the controller.
func (c *Controller) Stop() error {
	c.Manager.RemoveAllAndWait()
	return c.storeInformer.Stop()
}

// Name returns the name of the agent controller, it is completely defined by the name of
// the underlying store controller.
func (c *Controller) Name() string {
	return c.name
}

// newTriggerStoreInformer returns the trigger store controller for the agent.
// This controller watches for trigger objects in the store and performs action based
// on the changes to the object.
func (c *Controller) newTriggerStoreInformer() *store.Informer {
	return store.NewInformerWithSharedCache(
		fmt.Sprintf("%s/", v1alpha1.TriggerKeyPrefix),

		func(kv *v1alpha1.KVPairStruct) error {
			triggerName := strings.TrimPrefix(kv.Key, v1alpha1.TriggerKeyPrefix+"/")
			log.Infof("trigger added: %s", triggerName)
			c.mu.Lock()
			c.events[triggerName] = storeEvent{
				event: eventUpdate,
				kv:    kv,
			}
			c.mu.Unlock()
			return nil
		},

		// Delete function for the new agent.
		func(key string) error {
			triggerName := strings.TrimPrefix(key, v1alpha1.TriggerKeyPrefix+"/")
			log.Infof("trigger deleted: %s", triggerName)

			c.mu.Lock()
			c.events[triggerName] = storeEvent{
				event: eventRemove,
				kv: &v1alpha1.KVPairStruct{
					Key: key,
				},
			}
			c.mu.Unlock()
			return nil
		},

		func(kv *v1alpha1.KVPairStruct, version uint64) error {
			triggerName := strings.TrimPrefix(kv.Key, v1alpha1.TriggerKeyPrefix+"/")
			log.Infof("trigger added: %s", triggerName)
			c.mu.Lock()
			c.events[triggerName] = storeEvent{
				event: eventUpdate,
				kv:    kv,
			}
			c.mu.Unlock()
			return nil
		},
	)
}

func (c *Controller) addTrigger(kv *v1alpha1.KVPairStruct) error {
	triggerName := strings.TrimPrefix(kv.Key, v1alpha1.TriggerKeyPrefix+"/")

	var trigger v1alpha1.TriggerSpecWithName
	err := json.Unmarshal([]byte(kv.Value), &trigger)
	if err != nil {
		return fmt.Errorf("error while unmarshaling the trigger spec from data: %s", err)
	}

	if _, ok := c.Triggers[triggerName]; ok {
		err := c.Manager.RemoveController(triggerName)
		if err != nil {
			return fmt.Errorf("error while removing controller for %s: %s", triggerName, err)
		}
		delete(c.Triggers, triggerName)

		return nil
	}

	if _, ok := c.DefaultTriggers[triggerName]; ok {
		delete(c.DefaultTriggers, triggerName)
	}

	if trigger.Type == string(v1alpha1.DefaultTriggerType) {
		c.DefaultTriggers[triggerName] = struct{}{}
	} else {
		// Udpate the trigger specification
		t := Trigger{
			RunningPipelines:    0,
			UpdateAvailable:     false,
			TriggerSpecWithName: &trigger,
		}

		c.Triggers[triggerName] = &t
		err = t.SetupController(c.Manager)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) triggerUpdateControllerFunc(ctx context.Context) error {
	// Loop through each of the event and process them.
	for tName, ev := range c.events {
		switch ev.event {
		case eventUpdate:
			if trigger, ok := c.Triggers[tName]; ok {
				trigger.UpdateAvailable = true

				if trigger.RunningPipelines == 0 {
					err := c.addTrigger(ev.kv)
					if err != nil {
						log.Errorf("error while adding trigger(%s): %s", tName, err)
					} else {
						delete(c.events, tName)
					}
				}

				break
			}

			err := c.addTrigger(ev.kv)
			if err != nil {
				log.Errorf("error while adding trigger(%s): %s", tName, err)
			} else {
				delete(c.events, tName)
			}
		case eventRemove:
			if trigger, ok := c.Triggers[tName]; ok {
				trigger.UpdateAvailable = true

				if trigger.RunningPipelines == 0 {
					if err := trigger.StopController(c.Manager); err != nil {
						log.Errorf("%s", err)
					} else {
						delete(c.events, tName)
					}
				}

				break
			}

			delete(c.Triggers, tName)
			delete(c.DefaultTriggers, tName)
			delete(c.events, tName)
		default:
			return fmt.Errorf("Not a valid event(%s): %s", tName, ev.event)
		}
	}

	return nil
}
