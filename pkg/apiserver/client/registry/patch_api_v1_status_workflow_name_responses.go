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

// PatchAPIV1StatusWorkflowNameReader is a Reader for the PatchAPIV1StatusWorkflowName structure.
type PatchAPIV1StatusWorkflowNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchAPIV1StatusWorkflowNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 400:
		result := NewPatchAPIV1StatusWorkflowNameBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPatchAPIV1StatusWorkflowNameBadRequest creates a PatchAPIV1StatusWorkflowNameBadRequest with default headers values
func NewPatchAPIV1StatusWorkflowNameBadRequest() *PatchAPIV1StatusWorkflowNameBadRequest {
	return &PatchAPIV1StatusWorkflowNameBadRequest{}
}

/*PatchAPIV1StatusWorkflowNameBadRequest handles this case with default header values.

Bad Request
*/
type PatchAPIV1StatusWorkflowNameBadRequest struct {
	Payload *models.ResponseHTTPError
}

func (o *PatchAPIV1StatusWorkflowNameBadRequest) Error() string {
	return fmt.Sprintf("[PATCH /api/v1/status/workflow/{name}][%d] patchApiV1StatusWorkflowNameBadRequest  %+v", 400, o.Payload)
}

func (o *PatchAPIV1StatusWorkflowNameBadRequest) GetPayload() *models.ResponseHTTPError {
	return o.Payload
}

func (o *PatchAPIV1StatusWorkflowNameBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseHTTPError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
