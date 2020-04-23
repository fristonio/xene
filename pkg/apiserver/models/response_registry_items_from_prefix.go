// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ResponseRegistryItemsFromPrefix response registry items from prefix
//
// swagger:model response.RegistryItemsFromPrefix
type ResponseRegistryItemsFromPrefix struct {

	// count
	Count int64 `json:"count,omitempty"`

	// items
	Items string `json:"items,omitempty"`
}

// Validate validates this response registry items from prefix
func (m *ResponseRegistryItemsFromPrefix) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ResponseRegistryItemsFromPrefix) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseRegistryItemsFromPrefix) UnmarshalBinary(b []byte) error {
	var res ResponseRegistryItemsFromPrefix
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}