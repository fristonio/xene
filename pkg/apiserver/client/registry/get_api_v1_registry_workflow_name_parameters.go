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

// NewGetAPIV1RegistryWorkflowNameParams creates a new GetAPIV1RegistryWorkflowNameParams object
// with the default values initialized.
func NewGetAPIV1RegistryWorkflowNameParams() *GetAPIV1RegistryWorkflowNameParams {
	var ()
	return &GetAPIV1RegistryWorkflowNameParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAPIV1RegistryWorkflowNameParamsWithTimeout creates a new GetAPIV1RegistryWorkflowNameParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAPIV1RegistryWorkflowNameParamsWithTimeout(timeout time.Duration) *GetAPIV1RegistryWorkflowNameParams {
	var ()
	return &GetAPIV1RegistryWorkflowNameParams{

		timeout: timeout,
	}
}

// NewGetAPIV1RegistryWorkflowNameParamsWithContext creates a new GetAPIV1RegistryWorkflowNameParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAPIV1RegistryWorkflowNameParamsWithContext(ctx context.Context) *GetAPIV1RegistryWorkflowNameParams {
	var ()
	return &GetAPIV1RegistryWorkflowNameParams{

		Context: ctx,
	}
}

// NewGetAPIV1RegistryWorkflowNameParamsWithHTTPClient creates a new GetAPIV1RegistryWorkflowNameParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAPIV1RegistryWorkflowNameParamsWithHTTPClient(client *http.Client) *GetAPIV1RegistryWorkflowNameParams {
	var ()
	return &GetAPIV1RegistryWorkflowNameParams{
		HTTPClient: client,
	}
}

/*GetAPIV1RegistryWorkflowNameParams contains all the parameters to send to the API endpoint
for the get API v1 registry workflow name operation typically these are written to a http.Request
*/
type GetAPIV1RegistryWorkflowNameParams struct {

	/*Name
	  name of the workflow to get.

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) WithTimeout(timeout time.Duration) *GetAPIV1RegistryWorkflowNameParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) WithContext(ctx context.Context) *GetAPIV1RegistryWorkflowNameParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) WithHTTPClient(client *http.Client) *GetAPIV1RegistryWorkflowNameParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) WithName(name string) *GetAPIV1RegistryWorkflowNameParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get API v1 registry workflow name params
func (o *GetAPIV1RegistryWorkflowNameParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *GetAPIV1RegistryWorkflowNameParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
