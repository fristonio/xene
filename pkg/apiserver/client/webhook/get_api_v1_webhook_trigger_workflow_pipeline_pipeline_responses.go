// Code generated by go-swagger; DO NOT EDIT.

package webhook

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/fristonio/xene/pkg/apiserver/models"
)

// GetAPIV1WebhookTriggerWorkflowPipelinePipelineReader is a Reader for the GetAPIV1WebhookTriggerWorkflowPipelinePipeline structure.
type GetAPIV1WebhookTriggerWorkflowPipelinePipelineReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAPIV1WebhookTriggerWorkflowPipelinePipelineReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 400:
		result := NewGetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest creates a GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest with default headers values
func NewGetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest() *GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest {
	return &GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest{}
}

/*GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest handles this case with default header values.

Bad Request
*/
type GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/v1/webhook/trigger/{workflow}/pipeline/{pipeline}][%d] getApiV1WebhookTriggerWorkflowPipelinePipelineBadRequest  %+v", 400, o.Payload)
}

func (o *GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1WebhookTriggerWorkflowPipelinePipelineBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError creates a GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError with default headers values
func NewGetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError() *GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError {
	return &GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError{}
}

/*GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1/webhook/trigger/{workflow}/pipeline/{pipeline}][%d] getApiV1WebhookTriggerWorkflowPipelinePipelineInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1WebhookTriggerWorkflowPipelinePipelineInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
