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

// ResponsePipelineInfo response pipeline info
//
// swagger:model response.PipelineInfo
type ResponsePipelineInfo struct {

	// name
	Name string `json:"name,omitempty"`

	// runs
	Runs []*ResponsePipelineRunInfo `json:"runs"`

	// spec
	Spec string `json:"spec,omitempty"`

	// warnings
	Warnings []string `json:"warnings"`

	// workflow
	Workflow string `json:"workflow,omitempty"`
}

// Validate validates this response pipeline info
func (m *ResponsePipelineInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRuns(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ResponsePipelineInfo) validateRuns(formats strfmt.Registry) error {

	if swag.IsZero(m.Runs) { // not required
		return nil
	}

	for i := 0; i < len(m.Runs); i++ {
		if swag.IsZero(m.Runs[i]) { // not required
			continue
		}

		if m.Runs[i] != nil {
			if err := m.Runs[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("runs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ResponsePipelineInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponsePipelineInfo) UnmarshalBinary(b []byte) error {
	var res ResponsePipelineInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
