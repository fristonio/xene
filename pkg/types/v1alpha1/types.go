package v1alpha1

import "fmt"

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

// Validate checks for the information in TypeMeta
func (t *TypeMeta) Validate(reqKind string) error {
	if reqKind != t.Kind {
		return fmt.Errorf("kind is not valid for the type: %s(required %s)", t.Kind, reqKind)
	}

	return nil
}

// DeepEquals checks if the two TypeMeta objects are equal or not
func (t *TypeMeta) DeepEquals(tz *TypeMeta) bool {
	if t.Kind != tz.Kind ||
		t.APIVersion != tz.APIVersion {
		return false
	}

	return true
}

// ObjectMeta contains metadata which is common to all of the objects in xene.
type ObjectMeta struct {
	// Description for the type object.
	Description string `json:"description"`

	// Name is the name of the type object.
	Name string `json:"name"`
}

// GetName return the name of the object
func (o *ObjectMeta) GetName() string {
	return o.Name
}

// DeepEquals checks if the two ObjectMeta objects are equal or not
func (o *ObjectMeta) DeepEquals(oz *ObjectMeta) bool {
	if o.Name != oz.Name ||
		o.Description != o.Description {
		return false
	}

	return true
}

// Metadata corresponds to the metadata associated with the object.
type Metadata struct {
	ObjectMeta `json:",inline"`

	// UID is the unique ID associated with each managed resource.
	UID string `json:"uid,omitempty"`
}

// Validate validates the information present in metadata.
func (m *Metadata) Validate() error {
	if m.ObjectMeta.Name == "" {
		return fmt.Errorf("name is a required parameter in object manifest")
	}

	return nil
}

// RegistryItem is the type of registry item in xene.
type RegistryItem string
