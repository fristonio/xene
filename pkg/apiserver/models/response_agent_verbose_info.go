// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ResponseAgentVerboseInfo response agent verbose info
//
// swagger:model response.AgentVerboseInfo
type ResponseAgentVerboseInfo struct {

	// address
	Address string `json:"address,omitempty"`

	// healthy
	Healthy bool `json:"healthy,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// secrets
	Secrets []*ResponseAgentSecretInfo `json:"secrets"`

	// secure
	Secure bool `json:"secure,omitempty"`

	// server name
	ServerName string `json:"serverName,omitempty"`

	// workflows
	Workflows []*ResponseAgentWorkflowInfo `json:"workflows"`
}

// Validate validates this response agent verbose info
func (m *ResponseAgentVerboseInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSecrets(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWorkflows(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ResponseAgentVerboseInfo) validateSecrets(formats strfmt.Registry) error {

	if swag.IsZero(m.Secrets) { // not required
		return nil
	}

	for i := 0; i < len(m.Secrets); i++ {
		if swag.IsZero(m.Secrets[i]) { // not required
			continue
		}

		if m.Secrets[i] != nil {
			if err := m.Secrets[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("secrets" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ResponseAgentVerboseInfo) validateWorkflows(formats strfmt.Registry) error {

	if swag.IsZero(m.Workflows) { // not required
		return nil
	}

	for i := 0; i < len(m.Workflows); i++ {
		if swag.IsZero(m.Workflows[i]) { // not required
			continue
		}

		if m.Workflows[i] != nil {
			if err := m.Workflows[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("workflows" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ResponseAgentVerboseInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseAgentVerboseInfo) UnmarshalBinary(b []byte) error {
	var res ResponseAgentVerboseInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
