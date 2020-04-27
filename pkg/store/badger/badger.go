package badger

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/fristonio/xene/pkg/defaults"
	types "github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/fristonio/xene/pkg/utils"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrKeyAlreadyExist represent the error when the key is already present in the key
	// value store.
	ErrKeyAlreadyExist = errors.New("key already exist in the key value store")
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
func (b *Backend) Get(ctx context.Context, key string) (*types.Value, error) {
	txn := b.db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	val, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	return &types.Value{
		Data:             val,
		Version:          item.Version(),
		ExpiresAt:        item.ExpiresAt(),
		DeletedOrExpired: item.IsDeletedOrExpired(),
	}, nil
}

// Exists checks if the key provided exists in the datastore or not.
func (b *Backend) Exists(ctx context.Context, key string) (bool, error) {
	txn := b.db.NewTransaction(false)
	defer txn.Discard()

	_, err := txn.Get([]byte(key))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// KeyDoesNotExistError checks if the provided error is due to non existance of the key
// or not.
func (b *Backend) KeyDoesNotExistError(err error) bool {
	return err.Error() == badger.ErrKeyNotFound.Error()
}

// GetPrefix returns the first key which matches the prefix and its value
func (b *Backend) GetPrefix(ctx context.Context, path string) (string, *types.Value, error) {
	txn := b.db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(path)
	if it.Seek(prefix); it.ValidForPrefix(prefix) {
		item := it.Item()
		val, err := item.ValueCopy(nil)
		if err != nil {
			return "", nil, err
		}

		return string(item.Key()), &types.Value{
			Data:             val,
			Version:          item.Version(),
			ExpiresAt:        item.ExpiresAt(),
			DeletedOrExpired: item.IsDeletedOrExpired(),
		}, nil
	}

	return "", nil, nil
}

// Set sets value of key
func (b *Backend) Set(ctx context.Context, key string, value []byte) error {
	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	err := txn.Set([]byte(key), value)
	if err != nil {
		return err
	}

	return txn.Commit()
}

// Delete deletes a key
func (b *Backend) Delete(ctx context.Context, key string) error {
	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Delete([]byte(key)); err != nil {
		return err
	}

	return txn.Commit()
}

// DeletePrefix deletes the first key which matches the prefix and its value.
func (b *Backend) DeletePrefix(ctx context.Context, path string) error {
	txn := b.db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(path)
	if it.Seek(prefix); it.ValidForPrefix(prefix) {
		item := it.Item()
		t := b.db.NewTransaction(true)
		defer t.Discard()

		if err := t.Delete(item.Key()); err != nil {
			return err
		}

		return t.Commit()
	}

	return fmt.Errorf("no key with prefix %s found", path)
}

// CreateOnly atomically creates a key or fails if it already exists
func (b *Backend) CreateOnly(ctx context.Context, key string, value []byte) (bool, error) {
	_, err := b.Get(ctx, key)
	if err == nil {
		return true, ErrKeyAlreadyExist
	}

	err = b.Set(ctx, key, value)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateIfExists creates a key with the value only if key condKey exists
func (b *Backend) CreateIfExists(ctx context.Context, condKey, key string, value []byte) error {
	_, err := b.Get(ctx, condKey)
	if err != nil {
		return err
	}

	return b.Set(ctx, key, value)
}

// ListPrefix returns a list of keys matching the prefix
// This is best effort.
func (b *Backend) ListPrefix(ctx context.Context, path string) (types.KeyValuePairs, error) {
	list := map[string]types.Value{}
	txn := b.db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(path)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		val, err := item.ValueCopy(nil)
		if err != nil {
			log.Errorf("error while copying value for item: %s: %s", item.Key(), err)
			continue
		}

		list[string(item.Key())] = types.Value{
			Data:             val,
			Version:          item.Version(),
			ExpiresAt:        item.ExpiresAt(),
			DeletedOrExpired: item.IsDeletedOrExpired(),
		}
	}

	return list, nil
}

// PrefixScanWithFunction scans the provided key prefix keys in the store and run
// the provided function for each one of them.
func (b *Backend) PrefixScanWithFunction(ctx context.Context,
	key string, f types.KVPairStructFunc) {

	// Start a badger iteration and iterate over the target prefix keys.
	_ = b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(key)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				log.Errorf("error while copying value for item: %s: %s", key, err)
				continue
			}

			f(&types.KVPairStruct{
				Key:   string(item.Key()),
				Value: string(val),

				Version:          item.Version(),
				ExpiresAt:        item.ExpiresAt(),
				DeletedOrExpired: item.IsDeletedOrExpired(),
			})
		}
		return nil
	})
}

// Encode encodes a binary slice into a character set that the backend
// supports
func (b *Backend) Encode(in []byte) string {
	return base64.URLEncoding.EncodeToString([]byte(in))
}

// Decode decodes a key previously encoded back into the original binary slice
func (b *Backend) Decode(in string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(in)
}
