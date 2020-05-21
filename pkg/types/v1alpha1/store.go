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

// KVPairStruct is the representation of Key value pair in golang structure
// This is used when we serialize KeyValue to string.
type KVPairStruct struct {
	Key   string `json:"key" example:"registry/workflow/xxdfdihdfai=="`
	Value string `json:"value" example:"Workflow Document"`

	Version          uint64 `json:"version"`
	ExpiresAt        uint64 `json:"expiresAt"`
	DeletedOrExpired bool   `json:"deletedOrExpired"`
}

// KVPairStructFunc is the function type of a function which takes KVPairStruct as
// an argument.
type KVPairStructFunc func(*KVPairStruct)
