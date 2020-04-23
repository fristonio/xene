// Code generated by go-swagger; DO NOT EDIT.

package auth

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

// NewGetOauthProviderRedirectParams creates a new GetOauthProviderRedirectParams object
// with the default values initialized.
func NewGetOauthProviderRedirectParams() *GetOauthProviderRedirectParams {
	var ()
	return &GetOauthProviderRedirectParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetOauthProviderRedirectParamsWithTimeout creates a new GetOauthProviderRedirectParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetOauthProviderRedirectParamsWithTimeout(timeout time.Duration) *GetOauthProviderRedirectParams {
	var ()
	return &GetOauthProviderRedirectParams{

		timeout: timeout,
	}
}

// NewGetOauthProviderRedirectParamsWithContext creates a new GetOauthProviderRedirectParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetOauthProviderRedirectParamsWithContext(ctx context.Context) *GetOauthProviderRedirectParams {
	var ()
	return &GetOauthProviderRedirectParams{

		Context: ctx,
	}
}

// NewGetOauthProviderRedirectParamsWithHTTPClient creates a new GetOauthProviderRedirectParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetOauthProviderRedirectParamsWithHTTPClient(client *http.Client) *GetOauthProviderRedirectParams {
	var ()
	return &GetOauthProviderRedirectParams{
		HTTPClient: client,
	}
}

/*GetOauthProviderRedirectParams contains all the parameters to send to the API endpoint
for the get oauth provider redirect operation typically these are written to a http.Request
*/
type GetOauthProviderRedirectParams struct {

	/*Provider
	  Provider for the oauth redirect

	*/
	Provider string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) WithTimeout(timeout time.Duration) *GetOauthProviderRedirectParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) WithContext(ctx context.Context) *GetOauthProviderRedirectParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) WithHTTPClient(client *http.Client) *GetOauthProviderRedirectParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithProvider adds the provider to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) WithProvider(provider string) *GetOauthProviderRedirectParams {
	o.SetProvider(provider)
	return o
}

// SetProvider adds the provider to the get oauth provider redirect params
func (o *GetOauthProviderRedirectParams) SetProvider(provider string) {
	o.Provider = provider
}

// WriteToRequest writes these params to a swagger request
func (o *GetOauthProviderRedirectParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param provider
	qrProvider := o.Provider
	qProvider := qrProvider
	if qProvider != "" {
		if err := r.SetQueryParam("provider", qProvider); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
