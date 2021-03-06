// Code generated by go-swagger; DO NOT EDIT.

package status

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new status API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for status API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	DeleteAPIV1StatusWorkflowName(params *DeleteAPIV1StatusWorkflowNameParams, authInfo runtime.ClientAuthInfoWriter) (*DeleteAPIV1StatusWorkflowNameOK, error)

	GetAPIV1StatusWorkflow(params *GetAPIV1StatusWorkflowParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1StatusWorkflowOK, error)

	GetAPIV1StatusWorkflowName(params *GetAPIV1StatusWorkflowNameParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1StatusWorkflowNameOK, error)

	PostAPIV1StatusWorkflow(params *PostAPIV1StatusWorkflowParams, authInfo runtime.ClientAuthInfoWriter) (*PostAPIV1StatusWorkflowOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  DeleteAPIV1StatusWorkflowName deletes the specified workflow from the store

  Deletes the workflow status specified by the name parameter, if the workflow is not
*/
func (a *Client) DeleteAPIV1StatusWorkflowName(params *DeleteAPIV1StatusWorkflowNameParams, authInfo runtime.ClientAuthInfoWriter) (*DeleteAPIV1StatusWorkflowNameOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteAPIV1StatusWorkflowNameParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DeleteAPIV1StatusWorkflowName",
		Method:             "DELETE",
		PathPattern:        "/api/v1/status/workflow/{name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteAPIV1StatusWorkflowNameReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteAPIV1StatusWorkflowNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteAPIV1StatusWorkflowName: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetAPIV1StatusWorkflow returns the specified workflow status object from the store
*/
func (a *Client) GetAPIV1StatusWorkflow(params *GetAPIV1StatusWorkflowParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1StatusWorkflowOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetAPIV1StatusWorkflowParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetAPIV1StatusWorkflow",
		Method:             "GET",
		PathPattern:        "/api/v1/status/workflow",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetAPIV1StatusWorkflowReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetAPIV1StatusWorkflowOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetAPIV1StatusWorkflow: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetAPIV1StatusWorkflowName returns the specified workflow object from the store with the name in params
*/
func (a *Client) GetAPIV1StatusWorkflowName(params *GetAPIV1StatusWorkflowNameParams, authInfo runtime.ClientAuthInfoWriter) (*GetAPIV1StatusWorkflowNameOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetAPIV1StatusWorkflowNameParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetAPIV1StatusWorkflowName",
		Method:             "GET",
		PathPattern:        "/api/v1/status/workflow/{name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetAPIV1StatusWorkflowNameReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetAPIV1StatusWorkflowNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetAPIV1StatusWorkflowName: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  PostAPIV1StatusWorkflow creates a new workflow status in the store

  This route creates a new workflow status for xene to operate on, if the workflow already exists
*/
func (a *Client) PostAPIV1StatusWorkflow(params *PostAPIV1StatusWorkflowParams, authInfo runtime.ClientAuthInfoWriter) (*PostAPIV1StatusWorkflowOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostAPIV1StatusWorkflowParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "PostAPIV1StatusWorkflow",
		Method:             "POST",
		PathPattern:        "/api/v1/status/workflow",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/x-www-form-urlencoded"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PostAPIV1StatusWorkflowReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostAPIV1StatusWorkflowOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostAPIV1StatusWorkflow: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
