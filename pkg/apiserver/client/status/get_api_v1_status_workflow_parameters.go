// Code generated by go-swagger; DO NOT EDIT.

package status

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

// NewGetAPIV1StatusWorkflowParams creates a new GetAPIV1StatusWorkflowParams object
// with the default values initialized.
func NewGetAPIV1StatusWorkflowParams() *GetAPIV1StatusWorkflowParams {
	var ()
	return &GetAPIV1StatusWorkflowParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAPIV1StatusWorkflowParamsWithTimeout creates a new GetAPIV1StatusWorkflowParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAPIV1StatusWorkflowParamsWithTimeout(timeout time.Duration) *GetAPIV1StatusWorkflowParams {
	var ()
	return &GetAPIV1StatusWorkflowParams{

		timeout: timeout,
	}
}

// NewGetAPIV1StatusWorkflowParamsWithContext creates a new GetAPIV1StatusWorkflowParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAPIV1StatusWorkflowParamsWithContext(ctx context.Context) *GetAPIV1StatusWorkflowParams {
	var ()
	return &GetAPIV1StatusWorkflowParams{

		Context: ctx,
	}
}

// NewGetAPIV1StatusWorkflowParamsWithHTTPClient creates a new GetAPIV1StatusWorkflowParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAPIV1StatusWorkflowParamsWithHTTPClient(client *http.Client) *GetAPIV1StatusWorkflowParams {
	var ()
	return &GetAPIV1StatusWorkflowParams{
		HTTPClient: client,
	}
}

/*GetAPIV1StatusWorkflowParams contains all the parameters to send to the API endpoint
for the get API v1 status workflow operation typically these are written to a http.Request
*/
type GetAPIV1StatusWorkflowParams struct {

	/*Name
	  name of the workflow to get status object for.

	*/
	Name *string
	/*Prefix
	  Prefix based get for workflow status documents.

	*/
	Prefix *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) WithTimeout(timeout time.Duration) *GetAPIV1StatusWorkflowParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) WithContext(ctx context.Context) *GetAPIV1StatusWorkflowParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) WithHTTPClient(client *http.Client) *GetAPIV1StatusWorkflowParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) WithName(name *string) *GetAPIV1StatusWorkflowParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) SetName(name *string) {
	o.Name = name
}

// WithPrefix adds the prefix to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) WithPrefix(prefix *string) *GetAPIV1StatusWorkflowParams {
	o.SetPrefix(prefix)
	return o
}

// SetPrefix adds the prefix to the get API v1 status workflow params
func (o *GetAPIV1StatusWorkflowParams) SetPrefix(prefix *string) {
	o.Prefix = prefix
}

// WriteToRequest writes these params to a swagger request
func (o *GetAPIV1StatusWorkflowParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Name != nil {

		// query param name
		var qrName string
		if o.Name != nil {
			qrName = *o.Name
		}
		qName := qrName
		if qName != "" {
			if err := r.SetQueryParam("name", qName); err != nil {
				return err
			}
		}

	}

	if o.Prefix != nil {

		// query param prefix
		var qrPrefix string
		if o.Prefix != nil {
			qrPrefix = *o.Prefix
		}
		qPrefix := qrPrefix
		if qPrefix != "" {
			if err := r.SetQueryParam("prefix", qPrefix); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
