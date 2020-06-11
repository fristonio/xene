package v1alpha1

import "errors"

var (
	// ErrInvalidDNSSubdomainName is error used when the name is not a valid
	// DNS Subdomain conforming to RFC 1123
	ErrInvalidDNSSubdomainName = errors.New("Name should be a valid DNS subdomain conforming to RFC 1123")
)
