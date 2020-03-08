package store

import (
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/store/badger"
)

// KVStore is the configured backend store for xene.
var KVStore Backend

// Setup sets up the store backend for xene.
func Setup() error {
	if option.Config.Store.Engine == "" {
		return fmt.Errorf("no storage engine configured")
	}

	switch option.Config.Store.Engine {
	case defaults.StorageEngineBadger:
		backend, err := badger.NewBadgerBackend(option.Config.Store.StorageDirectory)
		if err != nil {
			return err
		}
		KVStore = backend
	default:
		return fmt.Errorf("not a supported storage engine: %s", option.Config.Store.Engine)
	}

	return nil
}
