package store

import (
	"fmt"
	"os"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/store/badger"
	"github.com/fristonio/xene/pkg/utils"
)

// KVStore is the configured backend store for xene.
var KVStore Backend

// Setup sets up the store backend for xene.
func Setup(storageDir string) error {
	if option.Config.Store.Engine == "" {
		return fmt.Errorf("no storage engine configured")
	}

	switch option.Config.Store.Engine {
	case defaults.StorageEngineBadger:
		if !utils.DirExists(storageDir) {
			if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
				return fmt.Errorf("error while creating storage directory: %s", err)
			}
		}
		backend, err := badger.NewBadgerBackend(storageDir)
		if err != nil {
			return err
		}
		KVStore = backend
	default:
		return fmt.Errorf("not a supported storage engine: %s", option.Config.Store.Engine)
	}

	return nil
}
