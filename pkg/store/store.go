package store

import (
	"context"

	"github.com/fristonio/xene/pkg/types/v1alpha1"
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
	Get(ctx context.Context, key string) ([]byte, error)

	// GetPrefix returns the first key which matches the prefix and its value
	GetPrefix(ctx context.Context, prefix string) (string, []byte, error)

	// Set sets value of key
	Set(ctx context.Context, key string, value []byte) error

	// Delete deletes a key
	Delete(ctx context.Context, key string) error

	// DeletePrefix deletes the first key which matches the prefix and its value.
	DeletePrefix(ctx context.Context, path string) error

	// Update atomically creates a key or fails if it already exists
	Update(ctx context.Context, key string, value []byte, lease bool) error

	// CreateOnly atomically creates a key or fails if it already exists
	CreateOnly(ctx context.Context, key string, value []byte, lease bool) (bool, error)

	// CreateIfExists creates a key with the value only if key condKey exists
	CreateIfExists(ctx context.Context, condKey, key string, value []byte, lease bool) error

	// ListPrefix returns a list of keys matching the prefix
	ListPrefix(ctx context.Context, prefix string) (v1alpha1.KeyValuePairs, error)

	// Encodes a binary slice into a character set that the backend
	// supports
	Encode(in []byte) string

	// Decodes a key previously encoded back into the original binary slice
	Decode(in string) ([]byte, error)
}
