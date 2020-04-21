package controller

import (
	"fmt"

	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"

	log "github.com/sirupsen/logrus"
)

// configureAgentController starts running the agent controller for apiserver.
// this controller watches for agent object in the store and perform action based
// on the changes to the object.
func configureAgentController() *store.Controller {
	return store.NewControllerWithSharedCache(
		fmt.Sprintf("%s/", v1alpha1.AgentKeyPrefix),
		true,
		func(kv *v1alpha1.KVPairStruct) {
			log.Infof("Agent add function: %v", kv)
		},
		func(kv *v1alpha1.KVPairStruct) {
			log.Infof("Agent delete function: %v", kv)
		},
		func(kv *v1alpha1.KVPairStruct, version uint64) {
			log.Infof("Agent update function: %v: previous version: %d", kv, version)
		},
	)
}

func init() {
	RegisteredControllers = append(RegisteredControllers, configureAgentController())
}
