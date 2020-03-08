package option

// Store is the type which store configuration data related to the storage
// data to be configured for xene.
type Store struct {
	// Engine is the storage engine to use for xene, it can be any of the previously configured store.
	Engine string `json:"engine"`

	// StorageDirectory is the directory to use for the storage engine configured
	// for xene.
	StorageDirectory string `json:"storageDirectory"`
}
