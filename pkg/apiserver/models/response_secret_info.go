// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ResponseSecretInfo response secret info
//
// swagger:model response.SecretInfo
type ResponseSecretInfo struct {

	// name
	Name string `json:"name,omitempty"`

	// restricted
	Restricted bool `json:"restricted,omitempty"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this response secret info
func (m *ResponseSecretInfo) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ResponseSecretInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseSecretInfo) UnmarshalBinary(b []byte) error {
	var res ResponseSecretInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
