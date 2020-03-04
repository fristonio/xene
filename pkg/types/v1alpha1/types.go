package v1alpha1

// TypeMeta describe the type details for a particular object understood by
// xene. It is same as that what kubernetes uses for type specification of the
// object.
type TypeMeta struct {
	// Kind is the kind of resource the object will represent.
	// For example this can be - Workflow
	Kind string `json:"kind,omitempty"`

	// ApiVersion is the version this specification corresponds to.
	APIVersion string `json:"apiVersion,omitempty"`
}

// Metadata corresponds to the metadata associated with the object.
type Metadata struct {
	// Name is the name of the resource created by xene
	Name string `json:"name,omitempty"`

	// UID is the unique ID associated with each managed resource.
	UID string `json:"uid,omitempty"`
}
