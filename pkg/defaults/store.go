package defaults

import "time"

const (
	// StorageDir is the directory to use for the local store to save
	// and data corresponding to xene.
	StorageDir string = "/var/run/xene/store"

	// StorageEngineBadger is the name of the stroage engine corresponding to
	// dgraph-io/badger key value store.
	StorageEngineBadger string = "badger"

	// StoreControllerRunInterval contains the run interval for store configured
	// controller do functions.
	StoreControllerRunInterval time.Duration = time.Second * 15

	// StoreControllerNameLength is the length of the store controller name
	StoreControllerNameLength uint32 = 12
)
