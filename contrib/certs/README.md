# Generate client and server certs

This makefile generates the root CA certificate, server certificate and client certificate which can be used
when configuring mutual TLS based authentication between component.

To generate the certificates run the following command:

```bash
$ make certs

Generating root CA
Generating server key and certificates
Generating client keys and certificates
```
