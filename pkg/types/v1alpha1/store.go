package v1alpha1

// Value is an abstraction of the data stored in the kvstore as well as the
// mod revision of that data.
type Value struct {
	// Data is data represented by the value.
	Data []byte

	// ModRevision is the revision of the modified key.
	ModRevision uint64

	// LeaseID is the ID of the lease associated with the key
	LeaseID int64
}

// KeyValuePairs is a map of key=value pairs
type KeyValuePairs map[string]Value
