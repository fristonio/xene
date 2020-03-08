package badger

import (
	"context"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/fristonio/xene/pkg/utils"
)

// NewBadgerBackend configures and returns badger backend.
func NewBadgerBackend(dir string) (*Backend, error) {
	if !utils.DirExists(dir) {
		return nil, fmt.Errorf("directory(%s) does not exists", dir)
	}
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		return nil, err
	}

	return &Backend{
		storageDir: dir,
		db:         db,
	}, nil
}

// Backend is the backend for badger key value store for xene.
type Backend struct {
	// storageDir is the storage directory for badger DB.
	storageDir string

	// db is the Database object for badger db.
	db *badger.DB
}

// GetName returns the name of the badger backend.
func (b *Backend) GetName() string {
	return defaults.StorageEngineBadger
}

// Close closes the database client and does not allow to do any more transactions.
func (b *Backend) Close() error {
	return b.db.Close()
}

// Configured returns if the backend client has been configured or not.
func (b *Backend) Configured() bool {
	return b.db != nil
}

// Status returns the status of the initialized backend, it returns an error
// if anything is wrong with the backend module.
func (b *Backend) Status() (string, error) {
	if b.db == nil {
		return fmt.Sprintf("badger database is configured"), nil
	}

	return "", fmt.Errorf("badger database is not configured")
}

// Get returns the value of the key.
func (b *Backend) Get(ctx context.Context, key string) ([]byte, error) {
	return []byte{}, nil
}

// GetPrefix returns the first key which matches the prefix and its value
func (b *Backend) GetPrefix(ctx context.Context, prefix string) (string, []byte, error) {
	return "", []byte{}, nil
}

// Set sets value of key
func (b *Backend) Set(ctx context.Context, key string, value []byte) error {
	return nil
}

// Delete deletes a key
func (b *Backend) Delete(ctx context.Context, key string) error {
	return nil
}

// DeletePrefix deletes the first key which matches the prefix and its value.
func (b *Backend) DeletePrefix(ctx context.Context, path string) error {
	return nil
}

// Update atomically creates a key or fails if it already exists
func (b *Backend) Update(ctx context.Context, key string, value []byte, lease bool) error {
	return nil
}

// CreateOnly atomically creates a key or fails if it already exists
func (b *Backend) CreateOnly(ctx context.Context, key string, value []byte, lease bool) (bool, error) {
	return false, nil
}

// CreateIfExists creates a key with the value only if key condKey exists
func (b *Backend) CreateIfExists(ctx context.Context, condKey, key string, value []byte, lease bool) error {
	return nil
}

// ListPrefix returns a list of keys matching the prefix
func (b *Backend) ListPrefix(ctx context.Context, prefix string) (v1alpha1.KeyValuePairs, error) {
	return v1alpha1.KeyValuePairs{}, nil
}

// Encode encodes a binary slice into a character set that the backend
// supports
func (b *Backend) Encode(in []byte) string {
	return ""
}

// Decode decodes a key previously encoded back into the original binary slice
func (b *Backend) Decode(in string) ([]byte, error) {
	return []byte{}, nil
}
