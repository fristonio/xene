package v1alpha1

import (
	"encoding/base64"
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
)

// Secret is the type which contains xene secret object definition.
type Secret struct {
	// TypeMeta stores meta type information for the agent object.
	TypeMeta `json:",inline"`

	// Metadata contains metadata about the Agent object.
	Metadata Metadata `json:"metadata"`

	// Spec contains the spec of the agent.
	Spec SecretSpec `json:"spec"`
}

// Validate validates the integrity of the Secret object.
func (s *Secret) Validate() error {
	err := s.TypeMeta.Validate(SecretKind)
	if err != nil {
		return err
	}

	err = s.Metadata.Validate()
	if err != nil {
		return err
	}

	return s.Spec.Validate()
}

// AddContent adds the provided data to the secret content after base64 encoding
// it to a string.
func (s *Secret) AddContent(data []byte) {
	s.Spec.Content = base64.StdEncoding.EncodeToString(data)
}

// GetContent adds the provided data to the secret content after base64 encoding
// it to a string.
func (s *Secret) GetContent() ([]byte, error) {
	return base64.StdEncoding.DecodeString(s.Spec.Content)
}

// SecretSpec contains the spec of the secret.
type SecretSpec struct {
	// Type contains the type of secret we are storing.
	Type string `json:"type"`

	// Content contains the base64 string representation of the secret
	// content.
	Content string `json:"content"`

	// Restricted returns if the secret is restricted or not.
	Restricted bool `json:"restricted"`
}

// Validate validates the integrity of the Secret object.
func (s *SecretSpec) Validate() error {
	if s.Type != defaults.SecretTypeCertificates && s.Type != defaults.SecretTypeDefault {
		return fmt.Errorf("Invalid %s secret type in the spec", s.Type)
	}

	_, err := base64.StdEncoding.DecodeString(s.Content)
	if err != nil {
		return fmt.Errorf("content does not have a valid base64 string")
	}

	return nil
}
