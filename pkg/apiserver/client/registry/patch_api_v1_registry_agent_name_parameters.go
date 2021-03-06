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

// NewPatchAPIV1RegistryAgentNameParams creates a new PatchAPIV1RegistryAgentNameParams object
// with the default values initialized.
func NewPatchAPIV1RegistryAgentNameParams() *PatchAPIV1RegistryAgentNameParams {
	var ()
	return &PatchAPIV1RegistryAgentNameParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPatchAPIV1RegistryAgentNameParamsWithTimeout creates a new PatchAPIV1RegistryAgentNameParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPatchAPIV1RegistryAgentNameParamsWithTimeout(timeout time.Duration) *PatchAPIV1RegistryAgentNameParams {
	var ()
	return &PatchAPIV1RegistryAgentNameParams{

		timeout: timeout,
	}
}

// NewPatchAPIV1RegistryAgentNameParamsWithContext creates a new PatchAPIV1RegistryAgentNameParams object
// with the default values initialized, and the ability to set a context for a request
func NewPatchAPIV1RegistryAgentNameParamsWithContext(ctx context.Context) *PatchAPIV1RegistryAgentNameParams {
	var ()
	return &PatchAPIV1RegistryAgentNameParams{

		Context: ctx,
	}
}

// NewPatchAPIV1RegistryAgentNameParamsWithHTTPClient creates a new PatchAPIV1RegistryAgentNameParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPatchAPIV1RegistryAgentNameParamsWithHTTPClient(client *http.Client) *PatchAPIV1RegistryAgentNameParams {
	var ()
	return &PatchAPIV1RegistryAgentNameParams{
		HTTPClient: client,
	}
}

/*PatchAPIV1RegistryAgentNameParams contains all the parameters to send to the API endpoint
for the patch API v1 registry agent name operation typically these are written to a http.Request
*/
type PatchAPIV1RegistryAgentNameParams struct {

	/*Name
	  Name of the agent to be patched.

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) WithTimeout(timeout time.Duration) *PatchAPIV1RegistryAgentNameParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) WithContext(ctx context.Context) *PatchAPIV1RegistryAgentNameParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) WithHTTPClient(client *http.Client) *PatchAPIV1RegistryAgentNameParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) WithName(name string) *PatchAPIV1RegistryAgentNameParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the patch API v1 registry agent name params
func (o *PatchAPIV1RegistryAgentNameParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *PatchAPIV1RegistryAgentNameParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
