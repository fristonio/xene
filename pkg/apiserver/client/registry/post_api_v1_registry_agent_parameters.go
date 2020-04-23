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

// NewPostAPIV1RegistryAgentParams creates a new PostAPIV1RegistryAgentParams object
// with the default values initialized.
func NewPostAPIV1RegistryAgentParams() *PostAPIV1RegistryAgentParams {
	var ()
	return &PostAPIV1RegistryAgentParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostAPIV1RegistryAgentParamsWithTimeout creates a new PostAPIV1RegistryAgentParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostAPIV1RegistryAgentParamsWithTimeout(timeout time.Duration) *PostAPIV1RegistryAgentParams {
	var ()
	return &PostAPIV1RegistryAgentParams{

		timeout: timeout,
	}
}

// NewPostAPIV1RegistryAgentParamsWithContext creates a new PostAPIV1RegistryAgentParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostAPIV1RegistryAgentParamsWithContext(ctx context.Context) *PostAPIV1RegistryAgentParams {
	var ()
	return &PostAPIV1RegistryAgentParams{

		Context: ctx,
	}
}

// NewPostAPIV1RegistryAgentParamsWithHTTPClient creates a new PostAPIV1RegistryAgentParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostAPIV1RegistryAgentParamsWithHTTPClient(client *http.Client) *PostAPIV1RegistryAgentParams {
	var ()
	return &PostAPIV1RegistryAgentParams{
		HTTPClient: client,
	}
}

/*PostAPIV1RegistryAgentParams contains all the parameters to send to the API endpoint
for the post API v1 registry agent operation typically these are written to a http.Request
*/
type PostAPIV1RegistryAgentParams struct {

	/*Agent
	  Agent manifest to be created.

	*/
	Agent string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) WithTimeout(timeout time.Duration) *PostAPIV1RegistryAgentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) WithContext(ctx context.Context) *PostAPIV1RegistryAgentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) WithHTTPClient(client *http.Client) *PostAPIV1RegistryAgentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAgent adds the agent to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) WithAgent(agent string) *PostAPIV1RegistryAgentParams {
	o.SetAgent(agent)
	return o
}

// SetAgent adds the agent to the post API v1 registry agent params
func (o *PostAPIV1RegistryAgentParams) SetAgent(agent string) {
	o.Agent = agent
}

// WriteToRequest writes these params to a swagger request
func (o *PostAPIV1RegistryAgentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// form param agent
	frAgent := o.Agent
	fAgent := frAgent
	if fAgent != "" {
		if err := r.SetFormParam("agent", fAgent); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
