// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ResponseAgentInfo response agent info
//
// swagger:model response.AgentInfo
type ResponseAgentInfo struct {

	// address
	Address string `json:"address,omitempty"`

	// available
	Available bool `json:"available,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// secure
	Secure bool `json:"secure,omitempty"`
}

// Validate validates this response agent info
func (m *ResponseAgentInfo) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ResponseAgentInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseAgentInfo) UnmarshalBinary(b []byte) error {
	var res ResponseAgentInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}