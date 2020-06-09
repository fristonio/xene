// Code generated by go-swagger; DO NOT EDIT.

package client

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiops "github.com/fristonio/xene/pkg/apiserver/client/api"
	"github.com/fristonio/xene/pkg/apiserver/client/auth"
	"github.com/fristonio/xene/pkg/apiserver/client/health"
	"github.com/fristonio/xene/pkg/apiserver/client/info"
	"github.com/fristonio/xene/pkg/apiserver/client/registry"
	"github.com/fristonio/xene/pkg/apiserver/client/status"
	"github.com/fristonio/xene/pkg/apiserver/client/webhook"
)

// Default xene API server HTTP client.
var Default = NewHTTPClient(nil)

const (
	// DefaultHost is the default Host
	// found in Meta (info) section of spec file
	DefaultHost string = "localhost:6060"
	// DefaultBasePath is the default BasePath
	// found in Meta (info) section of spec file
	DefaultBasePath string = "/"
)

// DefaultSchemes are the default schemes found in Meta (info) section of spec file
var DefaultSchemes = []string{"http"}

// NewHTTPClient creates a new xene API server HTTP client.
func NewHTTPClient(formats strfmt.Registry) *XeneAPIServer {
	return NewHTTPClientWithConfig(formats, nil)
}

// NewHTTPClientWithConfig creates a new xene API server HTTP client,
// using a customizable transport config.
func NewHTTPClientWithConfig(formats strfmt.Registry, cfg *TransportConfig) *XeneAPIServer {
	// ensure nullable parameters have default
	if cfg == nil {
		cfg = DefaultTransportConfig()
	}

	// create transport and client
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	return New(transport, formats)
}

// New creates a new xene API server client
func New(transport runtime.ClientTransport, formats strfmt.Registry) *XeneAPIServer {
	// ensure nullable parameters have default
	if formats == nil {
		formats = strfmt.Default
	}

	cli := new(XeneAPIServer)
	cli.Transport = transport
	cli.API = apiops.New(transport, formats)
	cli.Auth = auth.New(transport, formats)
	cli.Health = health.New(transport, formats)
	cli.Info = info.New(transport, formats)
	cli.Registry = registry.New(transport, formats)
	cli.Status = status.New(transport, formats)
	cli.Webhook = webhook.New(transport, formats)
	return cli
}

// DefaultTransportConfig creates a TransportConfig with the
// default settings taken from the meta section of the spec file.
func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		Host:     DefaultHost,
		BasePath: DefaultBasePath,
		Schemes:  DefaultSchemes,
	}
}

// TransportConfig contains the transport related info,
// found in the meta section of the spec file.
type TransportConfig struct {
	Host     string
	BasePath string
	Schemes  []string
}

// WithHost overrides the default host,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithHost(host string) *TransportConfig {
	cfg.Host = host
	return cfg
}

// WithBasePath overrides the default basePath,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithBasePath(basePath string) *TransportConfig {
	cfg.BasePath = basePath
	return cfg
}

// WithSchemes overrides the default schemes,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithSchemes(schemes []string) *TransportConfig {
	cfg.Schemes = schemes
	return cfg
}

// XeneAPIServer is a client for xene API server
type XeneAPIServer struct {
	API apiops.ClientService

	Auth auth.ClientService

	Health health.ClientService

	Info info.ClientService

	Registry registry.ClientService

	Status status.ClientService

	Webhook webhook.ClientService

	Transport runtime.ClientTransport
}

// SetTransport changes the transport on the client and all its subresources
func (c *XeneAPIServer) SetTransport(transport runtime.ClientTransport) {
	c.Transport = transport
	c.API.SetTransport(transport)
	c.Auth.SetTransport(transport)
	c.Health.SetTransport(transport)
	c.Info.SetTransport(transport)
	c.Registry.SetTransport(transport)
	c.Status.SetTransport(transport)
	c.Webhook.SetTransport(transport)
}
