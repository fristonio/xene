package v1alpha1

// Value is an abstraction of the data stored in the kvstore as well as the
// mod revision of that data.
type Value struct {
	// Data is data represented by the value.
	Data []byte

	// Version is the revision of the modified key.
	Version uint64

	// ExpiresAt is the time at which key will expire.
	ExpiresAt uint64

	// DeletedOrExpired checks if the key associated with the value is deleted
	// or expired.
	DeletedOrExpired bool
}

// KeyValuePairs is a map of key=value pairs
type KeyValuePairs map[string]Value
