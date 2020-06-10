// Code generated by go-swagger; DO NOT EDIT.

package info

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/fristonio/xene/pkg/apiserver/models"
)

// GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecReader is a Reader for the GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpec structure.
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK creates a GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK with default headers values
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK() *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK {
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK{}
}

/*GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK handles this case with default header values.

OK
*/
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK struct {
	Payload *models.ResponseRegistryItem
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/info/workflow/{workflow}/pipeline/{pipeline}/spec][%d] getApiV1InfoWorkflowWorkflowPipelinePipelineSpecOK  %+v", 200, o.Payload)
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK) GetPayload() *models.ResponseRegistryItem {
	return o.Payload
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseRegistryItem)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest creates a GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest with default headers values
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest() *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest {
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest{}
}

/*GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest handles this case with default header values.

Bad Request
*/
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/v1/info/workflow/{workflow}/pipeline/{pipeline}/spec][%d] getApiV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest  %+v", 400, o.Payload)
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError creates a GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError with default headers values
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError() *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError {
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError{}
}

/*GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1/info/workflow/{workflow}/pipeline/{pipeline}/spec][%d] getApiV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineSpecInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}