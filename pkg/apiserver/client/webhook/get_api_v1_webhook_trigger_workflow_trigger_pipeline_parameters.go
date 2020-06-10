// Code generated by go-swagger; DO NOT EDIT.

package webhook

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

// NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParams creates a new GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams object
// with the default values initialized.
func NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParams() *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	var ()
	return &GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParamsWithTimeout creates a new GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParamsWithTimeout(timeout time.Duration) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	var ()
	return &GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams{

		timeout: timeout,
	}
}

// NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParamsWithContext creates a new GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParamsWithContext(ctx context.Context) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	var ()
	return &GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams{

		Context: ctx,
	}
}

// NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParamsWithHTTPClient creates a new GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAPIV1WebhookTriggerWorkflowTriggerPipelineParamsWithHTTPClient(client *http.Client) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	var ()
	return &GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams{
		HTTPClient: client,
	}
}

/*GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams contains all the parameters to send to the API endpoint
for the get API v1 webhook trigger workflow trigger pipeline operation typically these are written to a http.Request
*/
type GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams struct {

	/*Pipeline
	  Name of the pipeline to be triggered.

	*/
	Pipeline string
	/*Trigger
	  Name of the trigger associated with the pipeline

	*/
	Trigger string
	/*Workflow
	  Name of the workflow.

	*/
	Workflow string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) WithTimeout(timeout time.Duration) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) WithContext(ctx context.Context) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) WithHTTPClient(client *http.Client) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPipeline adds the pipeline to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) WithPipeline(pipeline string) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	o.SetPipeline(pipeline)
	return o
}

// SetPipeline adds the pipeline to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) SetPipeline(pipeline string) {
	o.Pipeline = pipeline
}

// WithTrigger adds the trigger to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) WithTrigger(trigger string) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	o.SetTrigger(trigger)
	return o
}

// SetTrigger adds the trigger to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) SetTrigger(trigger string) {
	o.Trigger = trigger
}

// WithWorkflow adds the workflow to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) WithWorkflow(workflow string) *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams {
	o.SetWorkflow(workflow)
	return o
}

// SetWorkflow adds the workflow to the get API v1 webhook trigger workflow trigger pipeline params
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) SetWorkflow(workflow string) {
	o.Workflow = workflow
}

// WriteToRequest writes these params to a swagger request
func (o *GetAPIV1WebhookTriggerWorkflowTriggerPipelineParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param pipeline
	if err := r.SetPathParam("pipeline", o.Pipeline); err != nil {
		return err
	}

	// path param trigger
	if err := r.SetPathParam("trigger", o.Trigger); err != nil {
		return err
	}

	// path param workflow
	if err := r.SetPathParam("workflow", o.Workflow); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
