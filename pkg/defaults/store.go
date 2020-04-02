package defaults

const (
	// StorageDir is the directory to use for the local store to save
	// and data corresponding to xene.
	StorageDir string = "/var/run/xene/store"

	// StorageEngineBadger is the name of the stroage engine corresponding to
	// dgraph-io/badger key value store.
	StorageEngineBadger string = "badger"
)