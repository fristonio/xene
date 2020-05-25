// Code generated by go-swagger; DO NOT EDIT.

package info

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

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams creates a new GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams object
// with the default values initialized.
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams() *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	var ()
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParamsWithTimeout creates a new GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParamsWithTimeout(timeout time.Duration) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	var ()
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams{

		timeout: timeout,
	}
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParamsWithContext creates a new GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParamsWithContext(ctx context.Context) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	var ()
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams{

		Context: ctx,
	}
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParamsWithHTTPClient creates a new GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParamsWithHTTPClient(client *http.Client) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	var ()
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams{
		HTTPClient: client,
	}
}

/*GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams contains all the parameters to send to the API endpoint
for the get API v1 info workflow workflow pipeline pipeline runs run ID operation typically these are written to a http.Request
*/
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams struct {

	/*Pipeline
	  Name of the pipeline to return the info about.

	*/
	Pipeline string
	/*RunID
	  RUN ID of the pipeline run.

	*/
	RunID string
	/*Workflow
	  Name of the workflow to get information about.

	*/
	Workflow string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) WithTimeout(timeout time.Duration) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) WithContext(ctx context.Context) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) WithHTTPClient(client *http.Client) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPipeline adds the pipeline to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) WithPipeline(pipeline string) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	o.SetPipeline(pipeline)
	return o
}

// SetPipeline adds the pipeline to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) SetPipeline(pipeline string) {
	o.Pipeline = pipeline
}

// WithRunID adds the runID to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) WithRunID(runID string) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	o.SetRunID(runID)
	return o
}

// SetRunID adds the runId to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) SetRunID(runID string) {
	o.RunID = runID
}

// WithWorkflow adds the workflow to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) WithWorkflow(workflow string) *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams {
	o.SetWorkflow(workflow)
	return o
}

// SetWorkflow adds the workflow to the get API v1 info workflow workflow pipeline pipeline runs run ID params
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) SetWorkflow(workflow string) {
	o.Workflow = workflow
}

// WriteToRequest writes these params to a swagger request
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param pipeline
	if err := r.SetPathParam("pipeline", o.Pipeline); err != nil {
		return err
	}

	// path param runID
	if err := r.SetPathParam("runID", o.RunID); err != nil {
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
