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

// GetAPIV1RegistryListRegitemReader is a Reader for the GetAPIV1RegistryListRegitem structure.
type GetAPIV1RegistryListRegitemReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAPIV1RegistryListRegitemReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAPIV1RegistryListRegitemOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetAPIV1RegistryListRegitemBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAPIV1RegistryListRegitemInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAPIV1RegistryListRegitemOK creates a GetAPIV1RegistryListRegitemOK with default headers values
func NewGetAPIV1RegistryListRegitemOK() *GetAPIV1RegistryListRegitemOK {
	return &GetAPIV1RegistryListRegitemOK{}
}

/*GetAPIV1RegistryListRegitemOK handles this case with default header values.

OK
*/
type GetAPIV1RegistryListRegitemOK struct {
	Payload []*models.ResponseAgentInfo
}

func (o *GetAPIV1RegistryListRegitemOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/registry/list/{regitem}][%d] getApiV1RegistryListRegitemOK  %+v", 200, o.Payload)
}

func (o *GetAPIV1RegistryListRegitemOK) GetPayload() []*models.ResponseAgentInfo {
	return o.Payload
}

func (o *GetAPIV1RegistryListRegitemOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1RegistryListRegitemBadRequest creates a GetAPIV1RegistryListRegitemBadRequest with default headers values
func NewGetAPIV1RegistryListRegitemBadRequest() *GetAPIV1RegistryListRegitemBadRequest {
	return &GetAPIV1RegistryListRegitemBadRequest{}
}

/*GetAPIV1RegistryListRegitemBadRequest handles this case with default header values.

Bad Request
*/
type GetAPIV1RegistryListRegitemBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1RegistryListRegitemBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/v1/registry/list/{regitem}][%d] getApiV1RegistryListRegitemBadRequest  %+v", 400, o.Payload)
}

func (o *GetAPIV1RegistryListRegitemBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1RegistryListRegitemBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAPIV1RegistryListRegitemInternalServerError creates a GetAPIV1RegistryListRegitemInternalServerError with default headers values
func NewGetAPIV1RegistryListRegitemInternalServerError() *GetAPIV1RegistryListRegitemInternalServerError {
	return &GetAPIV1RegistryListRegitemInternalServerError{}
}

/*GetAPIV1RegistryListRegitemInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetAPIV1RegistryListRegitemInternalServerError struct {
	Payload *models.ResponseHTTPError
}

func (o *GetAPIV1RegistryListRegitemInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1/registry/list/{regitem}][%d] getApiV1RegistryListRegitemInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAPIV1RegistryListRegitemInternalServerError) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *GetAPIV1RegistryListRegitemInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
