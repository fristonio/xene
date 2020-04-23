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

// DeleteAPIV1RegistryAgentNameReader is a Reader for the DeleteAPIV1RegistryAgentName structure.
type DeleteAPIV1RegistryAgentNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteAPIV1RegistryAgentNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteAPIV1RegistryAgentNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteAPIV1RegistryAgentNameBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteAPIV1RegistryAgentNameInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDeleteAPIV1RegistryAgentNameOK creates a DeleteAPIV1RegistryAgentNameOK with default headers values
func NewDeleteAPIV1RegistryAgentNameOK() *DeleteAPIV1RegistryAgentNameOK {
	return &DeleteAPIV1RegistryAgentNameOK{}
}

/*DeleteAPIV1RegistryAgentNameOK handles this case with default header values.

OK
*/
type DeleteAPIV1RegistryAgentNameOK struct {
	Payload *models.ResponseHTTPMessage
}

func (o *DeleteAPIV1RegistryAgentNameOK) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/registry/agent/{name}][%d] deleteApiV1RegistryAgentNameOK  %+v", 200, o.Payload)
}

func (o *DeleteAPIV1RegistryAgentNameOK) GetPayload() *models.ResponseHTTPMessage {
	return o.Payload
}

func (o *DeleteAPIV1RegistryAgentNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPMessage)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteAPIV1RegistryAgentNameBadRequest creates a DeleteAPIV1RegistryAgentNameBadRequest with default headers values
func NewDeleteAPIV1RegistryAgentNameBadRequest() *DeleteAPIV1RegistryAgentNameBadRequest {
	return &DeleteAPIV1RegistryAgentNameBadRequest{}
}

/*DeleteAPIV1RegistryAgentNameBadRequest handles this case with default header values.

Bad Request
*/
type DeleteAPIV1RegistryAgentNameBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *DeleteAPIV1RegistryAgentNameBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/registry/agent/{name}][%d] deleteApiV1RegistryAgentNameBadRequest  %+v", 400, o.Payload)
}

func (o *DeleteAPIV1RegistryAgentNameBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *DeleteAPIV1RegistryAgentNameBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteAPIV1RegistryAgentNameInternalServerError creates a DeleteAPIV1RegistryAgentNameInternalServerError with default headers values
func NewDeleteAPIV1RegistryAgentNameInternalServerError() *DeleteAPIV1RegistryAgentNameInternalServerError {
	return &DeleteAPIV1RegistryAgentNameInternalServerError{}
}

/*DeleteAPIV1RegistryAgentNameInternalServerError handles this case with default header values.

Internal Server Error
*/
type DeleteAPIV1RegistryAgentNameInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *DeleteAPIV1RegistryAgentNameInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/registry/agent/{name}][%d] deleteApiV1RegistryAgentNameInternalServerError  %+v", 500, o.Payload)
}

func (o *DeleteAPIV1RegistryAgentNameInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *DeleteAPIV1RegistryAgentNameInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
