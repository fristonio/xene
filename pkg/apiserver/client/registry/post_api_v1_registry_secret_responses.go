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

// PostAPIV1RegistrySecretReader is a Reader for the PostAPIV1RegistrySecret structure.
type PostAPIV1RegistrySecretReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostAPIV1RegistrySecretReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostAPIV1RegistrySecretOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostAPIV1RegistrySecretBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPostAPIV1RegistrySecretInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPostAPIV1RegistrySecretOK creates a PostAPIV1RegistrySecretOK with default headers values
func NewPostAPIV1RegistrySecretOK() *PostAPIV1RegistrySecretOK {
	return &PostAPIV1RegistrySecretOK{}
}

/*PostAPIV1RegistrySecretOK handles this case with default header values.

OK
*/
type PostAPIV1RegistrySecretOK struct {
	Payload *models.ResponseHTTPMessage
}

func (o *PostAPIV1RegistrySecretOK) Error() string {
	return fmt.Sprintf("[POST /api/v1/registry/secret][%d] postApiV1RegistrySecretOK  %+v", 200, o.Payload)
}

func (o *PostAPIV1RegistrySecretOK) GetPayload() *models.ResponseHTTPMessage {
	return o.Payload
}

func (o *PostAPIV1RegistrySecretOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPMessage)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostAPIV1RegistrySecretBadRequest creates a PostAPIV1RegistrySecretBadRequest with default headers values
func NewPostAPIV1RegistrySecretBadRequest() *PostAPIV1RegistrySecretBadRequest {
	return &PostAPIV1RegistrySecretBadRequest{}
}

/*PostAPIV1RegistrySecretBadRequest handles this case with default header values.

Bad Request
*/
type PostAPIV1RegistrySecretBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *PostAPIV1RegistrySecretBadRequest) Error() string {
	return fmt.Sprintf("[POST /api/v1/registry/secret][%d] postApiV1RegistrySecretBadRequest  %+v", 400, o.Payload)
}

func (o *PostAPIV1RegistrySecretBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *PostAPIV1RegistrySecretBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostAPIV1RegistrySecretInternalServerError creates a PostAPIV1RegistrySecretInternalServerError with default headers values
func NewPostAPIV1RegistrySecretInternalServerError() *PostAPIV1RegistrySecretInternalServerError {
	return &PostAPIV1RegistrySecretInternalServerError{}
}

/*PostAPIV1RegistrySecretInternalServerError handles this case with default header values.

Internal Server Error
*/
type PostAPIV1RegistrySecretInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *PostAPIV1RegistrySecretInternalServerError) Error() string {
	return fmt.Sprintf("[POST /api/v1/registry/secret][%d] postApiV1RegistrySecretInternalServerError  %+v", 500, o.Payload)
}

func (o *PostAPIV1RegistrySecretInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *PostAPIV1RegistrySecretInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
