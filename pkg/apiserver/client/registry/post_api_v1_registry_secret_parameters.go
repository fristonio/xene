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

// NewPostAPIV1RegistrySecretParams creates a new PostAPIV1RegistrySecretParams object
// with the default values initialized.
func NewPostAPIV1RegistrySecretParams() *PostAPIV1RegistrySecretParams {
	var ()
	return &PostAPIV1RegistrySecretParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostAPIV1RegistrySecretParamsWithTimeout creates a new PostAPIV1RegistrySecretParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostAPIV1RegistrySecretParamsWithTimeout(timeout time.Duration) *PostAPIV1RegistrySecretParams {
	var ()
	return &PostAPIV1RegistrySecretParams{

		timeout: timeout,
	}
}

// NewPostAPIV1RegistrySecretParamsWithContext creates a new PostAPIV1RegistrySecretParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostAPIV1RegistrySecretParamsWithContext(ctx context.Context) *PostAPIV1RegistrySecretParams {
	var ()
	return &PostAPIV1RegistrySecretParams{

		Context: ctx,
	}
}

// NewPostAPIV1RegistrySecretParamsWithHTTPClient creates a new PostAPIV1RegistrySecretParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostAPIV1RegistrySecretParamsWithHTTPClient(client *http.Client) *PostAPIV1RegistrySecretParams {
	var ()
	return &PostAPIV1RegistrySecretParams{
		HTTPClient: client,
	}
}

/*PostAPIV1RegistrySecretParams contains all the parameters to send to the API endpoint
for the post API v1 registry secret operation typically these are written to a http.Request
*/
type PostAPIV1RegistrySecretParams struct {

	/*Secret
	  secret manifest to be created.

	*/
	Secret string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) WithTimeout(timeout time.Duration) *PostAPIV1RegistrySecretParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) WithContext(ctx context.Context) *PostAPIV1RegistrySecretParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) WithHTTPClient(client *http.Client) *PostAPIV1RegistrySecretParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSecret adds the secret to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) WithSecret(secret string) *PostAPIV1RegistrySecretParams {
	o.SetSecret(secret)
	return o
}

// SetSecret adds the secret to the post API v1 registry secret params
func (o *PostAPIV1RegistrySecretParams) SetSecret(secret string) {
	o.Secret = secret
}

// WriteToRequest writes these params to a swagger request
func (o *PostAPIV1RegistrySecretParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// form param secret
	frSecret := o.Secret
	fSecret := frSecret
	if fSecret != "" {
		if err := r.SetFormParam("secret", fSecret); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
