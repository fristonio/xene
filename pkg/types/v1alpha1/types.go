package v1alpha1

// TypeMeta describe the type details for a particular object understood by
// xene. It is same as that what kubernetes uses for type specification of the
// object.
type TypeMeta struct {
	// Kind is the kind of resource the object will represent.
	// For example this can be - Workflow
	Kind string `json:"kind"`

	// ApiVersion is the version this specification corresponds to.
	APIVersion string `json:"apiVersion"`
}

// ObjectMeta contains metadata which is common to all of the objects in xene.
type ObjectMeta struct {
	// Description for the type object.
	Description string `json:"description"`

	// Name is the name of the type object.
	Name string `json:"name"`
}

// Metadata corresponds to the metadata associated with the object.
type Metadata struct {
	ObjectMeta `json:",inline"`

	// UID is the unique ID associated with each managed resource.
	UID string `json:"uid,omitempty"`
}
