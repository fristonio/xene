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

// GetAPIV1RegistryListAgentsReader is a Reader for the GetAPIV1RegistryListAgents structure.
type GetAPIV1RegistryListAgentsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAPIV1RegistryListAgentsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAPIV1RegistryListAgentsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAPIV1RegistryListAgentsOK creates a GetAPIV1RegistryListAgentsOK with default headers values
func NewGetAPIV1RegistryListAgentsOK() *GetAPIV1RegistryListAgentsOK {
	return &GetAPIV1RegistryListAgentsOK{}
}

/*GetAPIV1RegistryListAgentsOK handles this case with default header values.

OK
*/
type GetAPIV1RegistryListAgentsOK struct {
	Payload []*models.ResponseAgentInfo
}

func (o *GetAPIV1RegistryListAgentsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/registry/list/agents][%d] getApiV1RegistryListAgentsOK  %+v", 200, o.Payload)
}

func (o *GetAPIV1RegistryListAgentsOK) GetPayload() []*models.ResponseAgentInfo {
	return o.Payload
}

func (o *GetAPIV1RegistryListAgentsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
