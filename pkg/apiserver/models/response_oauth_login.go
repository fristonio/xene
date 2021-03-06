// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ResponseOauthLogin response oauth login
//
// swagger:model response.OauthLogin
type ResponseOauthLogin struct {

	// LoginURL is the URL to be used for logging in.
	LoginURL string `json:"loginURL,omitempty"`
}

// Validate validates this response oauth login
func (m *ResponseOauthLogin) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ResponseOauthLogin) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseOauthLogin) UnmarshalBinary(b []byte) error {
	var res ResponseOauthLogin
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
