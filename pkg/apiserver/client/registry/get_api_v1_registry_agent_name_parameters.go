// Code generated by go-swagger; DO NOT EDIT.

package registry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetAPIV1RegistryAgentNameParams creates a new GetAPIV1RegistryAgentNameParams object
// with the default values initialized.
func NewGetAPIV1RegistryAgentNameParams() *GetAPIV1RegistryAgentNameParams {
	var ()
	return &GetAPIV1RegistryAgentNameParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAPIV1RegistryAgentNameParamsWithTimeout creates a new GetAPIV1RegistryAgentNameParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAPIV1RegistryAgentNameParamsWithTimeout(timeout time.Duration) *GetAPIV1RegistryAgentNameParams {
	var ()
	return &GetAPIV1RegistryAgentNameParams{

		timeout: timeout,
	}
}

// NewGetAPIV1RegistryAgentNameParamsWithContext creates a new GetAPIV1RegistryAgentNameParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAPIV1RegistryAgentNameParamsWithContext(ctx context.Context) *GetAPIV1RegistryAgentNameParams {
	var ()
	return &GetAPIV1RegistryAgentNameParams{

		Context: ctx,
	}
}

// NewGetAPIV1RegistryAgentNameParamsWithHTTPClient creates a new GetAPIV1RegistryAgentNameParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAPIV1RegistryAgentNameParamsWithHTTPClient(client *http.Client) *GetAPIV1RegistryAgentNameParams {
	var ()
	return &GetAPIV1RegistryAgentNameParams{
		HTTPClient: client,
	}
}

/*GetAPIV1RegistryAgentNameParams contains all the parameters to send to the API endpoint
for the get API v1 registry agent name operation typically these are written to a http.Request
*/
type GetAPIV1RegistryAgentNameParams struct {

	/*Name
	  name of the agent to get.

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) WithTimeout(timeout time.Duration) *GetAPIV1RegistryAgentNameParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) WithContext(ctx context.Context) *GetAPIV1RegistryAgentNameParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) WithHTTPClient(client *http.Client) *GetAPIV1RegistryAgentNameParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) WithName(name string) *GetAPIV1RegistryAgentNameParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get API v1 registry agent name params
func (o *GetAPIV1RegistryAgentNameParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *GetAPIV1RegistryAgentNameParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
