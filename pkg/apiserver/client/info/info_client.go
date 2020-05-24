// Code generated by go-swagger; DO NOT EDIT.

package info

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new info API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for info API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	GetAPIV1InfoAgentName(params *GetAPIV1InfoAgentNameParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1InfoAgentNameOK, error)

	GetAPIV1InfoWorkflowNamePipelinePipeline(params *GetAPIV1InfoWorkflowNamePipelinePipelineParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1InfoWorkflowNamePipelinePipelineOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetAPIV1InfoAgentName returns verbose information about the agent
*/
func (a *Client) GetAPIV1InfoAgentName(params *GetAPIV1InfoAgentNameParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1InfoAgentNameOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetAPIV1InfoAgentNameParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetAPIV1InfoAgentName",
		Method:             "GET",
		PathPattern:        "/api/v1/info/agent/{name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetAPIV1InfoAgentNameReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetAPIV1InfoAgentNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetAPIV1InfoAgentName: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetAPIV1InfoWorkflowNamePipelinePipeline returns verbose information about a workflow
*/
func (a *Client) GetAPIV1InfoWorkflowNamePipelinePipeline(params *GetAPIV1InfoWorkflowNamePipelinePipelineParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1InfoWorkflowNamePipelinePipelineOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetAPIV1InfoWorkflowNamePipelinePipelineParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetAPIV1InfoWorkflowNamePipelinePipeline",
		Method:             "GET",
		PathPattern:        "/api/v1/info/workflow/{name}/pipeline/{pipeline}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetAPIV1InfoWorkflowNamePipelinePipelineReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetAPIV1InfoWorkflowNamePipelinePipelineOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetAPIV1InfoWorkflowNamePipelinePipeline: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
