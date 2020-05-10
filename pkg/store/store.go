package store

import (
	"context"

	types "github.com/fristonio/xene/pkg/types/v1alpha1"
)

// Backend is the interface which needs to be implemented by each client configured
// for the store.
type Backend interface {
	// GetName returns the name of the backend module.
	GetName() string

	// Close closes the database client and does not allow to do any more transactions.
	Close() error

	// Configured returns if the backend client has been configured or not.
	Configured() bool

	// Status returns the status of the initialized backend, it returns an error
	// if anything is wrong with the backend module.
	Status() (string, error)

	// Get returns the value of the key.
	Get(ctx context.Context, key string) (*types.Value, error)

	// Exists checks if the provided key exists or not
	Exists(ctx context.Context, key string) (bool, error)

	// KeyDoesNotExistError is the method to check if the error is due to non
	// existance of the key during the get operation.
	KeyDoesNotExistError(err error) bool

	// GetPrefix returns the first key which matches the prefix and its value
	GetPrefix(ctx context.Context, prefix string) (string, *types.Value, error)

	// Set sets value of key
	Set(ctx context.Context, key string, value []byte) error

	// Delete deletes a key
	Delete(ctx context.Context, key string) error

	// DeletePrefix deletes the first key which matches the prefix and its value.
	DeletePrefix(ctx context.Context, path string) error

	// CreateOnly atomically creates a key or fails if it already exists
	CreateOnly(ctx context.Context, key string, value []byte) (bool, error)

	// CreateIfExists creates a key with the value only if key condKey exists
	CreateIfExists(ctx context.Context, condKey, key string, value []byte) error

	// ListPrefixKeys list all the keys with the provided prefix.
	ListPrefixKeys(ctx context.Context, path string) ([]string, error)

	// ListPrefix returns a list of keys matching the prefix
	ListPrefix(ctx context.Context, prefix string) (types.KeyValuePairs, error)

	// ListPrefix returns a list of keys matching the prefix
	PrefixScanWithFunction(ctx context.Context, prefix string, f types.KVPairStructFunc)

	// Encodes a binary slice into a character set that the backend
	// supports
	Encode(in []byte) string

	// Decodes a key previously encoded back into the original binary slice
	Decode(in string) ([]byte, error)
}
