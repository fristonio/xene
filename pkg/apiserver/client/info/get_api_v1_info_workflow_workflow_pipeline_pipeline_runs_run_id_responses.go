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

// GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDReader is a Reader for the GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunID structure.
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK creates a GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK with default headers values
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK() *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK {
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK{}
}

/*GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK handles this case with default header values.

OK
*/
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK struct {
	Payload *models.ResponsePipelineRunVerboseInfo
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/info/workflow/{workflow}/pipeline/{pipeline}/runs/{runID}][%d] getApiV1InfoWorkflowWorkflowPipelinePipelineRunsRunIdOK  %+v", 200, o.Payload)
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK) GetPayload() *models.ResponsePipelineRunVerboseInfo {
	return o.Payload
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponsePipelineRunVerboseInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest creates a GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest with default headers values
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest() *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest {
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest{}
}

/*GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest handles this case with default header values.

Bad Request
*/
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/v1/info/workflow/{workflow}/pipeline/{pipeline}/runs/{runID}][%d] getApiV1InfoWorkflowWorkflowPipelinePipelineRunsRunIdBadRequest  %+v", 400, o.Payload)
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError creates a GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError with default headers values
func NewGetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError() *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError {
	return &GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError{}
}

/*GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1/info/workflow/{workflow}/pipeline/{pipeline}/runs/{runID}][%d] getApiV1InfoWorkflowWorkflowPipelinePipelineRunsRunIdInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1InfoWorkflowWorkflowPipelinePipelineRunsRunIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
