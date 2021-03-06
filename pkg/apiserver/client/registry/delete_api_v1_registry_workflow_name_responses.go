// Code generated by go-swagger; DO NOT EDIT.

package registry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/fristonio/xene/pkg/apiserver/models"
)

// DeleteAPIV1RegistryWorkflowNameReader is a Reader for the DeleteAPIV1RegistryWorkflowName structure.
type DeleteAPIV1RegistryWorkflowNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteAPIV1RegistryWorkflowNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteAPIV1RegistryWorkflowNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteAPIV1RegistryWorkflowNameBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteAPIV1RegistryWorkflowNameInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDeleteAPIV1RegistryWorkflowNameOK creates a DeleteAPIV1RegistryWorkflowNameOK with default headers values
func NewDeleteAPIV1RegistryWorkflowNameOK() *DeleteAPIV1RegistryWorkflowNameOK {
	return &DeleteAPIV1RegistryWorkflowNameOK{}
}

/*DeleteAPIV1RegistryWorkflowNameOK handles this case with default header values.

OK
*/
type DeleteAPIV1RegistryWorkflowNameOK struct {
	Payload *models.ResponseHTTPMessage
}

func (o *DeleteAPIV1RegistryWorkflowNameOK) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/registry/workflow/{name}][%d] deleteApiV1RegistryWorkflowNameOK  %+v", 200, o.Payload)
}

func (o *DeleteAPIV1RegistryWorkflowNameOK) GetPayload() *models.ResponseHTTPMessage {
	return o.Payload
}

func (o *DeleteAPIV1RegistryWorkflowNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPMessage)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteAPIV1RegistryWorkflowNameBadRequest creates a DeleteAPIV1RegistryWorkflowNameBadRequest with default headers values
func NewDeleteAPIV1RegistryWorkflowNameBadRequest() *DeleteAPIV1RegistryWorkflowNameBadRequest {
	return &DeleteAPIV1RegistryWorkflowNameBadRequest{}
}

/*DeleteAPIV1RegistryWorkflowNameBadRequest handles this case with default header values.

Bad Request
*/
type DeleteAPIV1RegistryWorkflowNameBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *DeleteAPIV1RegistryWorkflowNameBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/registry/workflow/{name}][%d] deleteApiV1RegistryWorkflowNameBadRequest  %+v", 400, o.Payload)
}

func (o *DeleteAPIV1RegistryWorkflowNameBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *DeleteAPIV1RegistryWorkflowNameBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteAPIV1RegistryWorkflowNameInternalServerError creates a DeleteAPIV1RegistryWorkflowNameInternalServerError with default headers values
func NewDeleteAPIV1RegistryWorkflowNameInternalServerError() *DeleteAPIV1RegistryWorkflowNameInternalServerError {
	return &DeleteAPIV1RegistryWorkflowNameInternalServerError{}
}

/*DeleteAPIV1RegistryWorkflowNameInternalServerError handles this case with default header values.

Internal Server Error
*/
type DeleteAPIV1RegistryWorkflowNameInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *DeleteAPIV1RegistryWorkflowNameInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/registry/workflow/{name}][%d] deleteApiV1RegistryWorkflowNameInternalServerError  %+v", 500, o.Payload)
}

func (o *DeleteAPIV1RegistryWorkflowNameInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *DeleteAPIV1RegistryWorkflowNameInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
