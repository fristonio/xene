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

// GetAPIV1RegistryAgentNameReader is a Reader for the GetAPIV1RegistryAgentName structure.
type GetAPIV1RegistryAgentNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAPIV1RegistryAgentNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAPIV1RegistryAgentNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetAPIV1RegistryAgentNameInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAPIV1RegistryAgentNameOK creates a GetAPIV1RegistryAgentNameOK with default headers values
func NewGetAPIV1RegistryAgentNameOK() *GetAPIV1RegistryAgentNameOK {
	return &GetAPIV1RegistryAgentNameOK{}
}

/*GetAPIV1RegistryAgentNameOK handles this case with default header values.

OK
*/
type GetAPIV1RegistryAgentNameOK struct {
	Payload *models.ResponseRegistryItem
}

func (o *GetAPIV1RegistryAgentNameOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/registry/agent/{name}][%d] getApiV1RegistryAgentNameOK  %+v", 200, o.Payload)
}

func (o *GetAPIV1RegistryAgentNameOK) GetPayload() *models.ResponseRegistryItem {
	return o.Payload
}

func (o *GetAPIV1RegistryAgentNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseRegistryItem)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1RegistryAgentNameInternalServerError creates a GetAPIV1RegistryAgentNameInternalServerError with default headers values
func NewGetAPIV1RegistryAgentNameInternalServerError() *GetAPIV1RegistryAgentNameInternalServerError {
	return &GetAPIV1RegistryAgentNameInternalServerError{}
}

/*GetAPIV1RegistryAgentNameInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetAPIV1RegistryAgentNameInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1RegistryAgentNameInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1/registry/agent/{name}][%d] getApiV1RegistryAgentNameInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAPIV1RegistryAgentNameInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1RegistryAgentNameInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
