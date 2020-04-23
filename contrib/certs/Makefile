SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

certs:
> @echo "Generating root CA"
> @mkdir -p ca.gen/ && \
	cd ca.gen/ && \
	cfssl gencert -initca ./../ca.json | cfssljson -bare ca && \
	cd ../
> @echo "Generating server key and certificates"
> @mkdir -p server.gen/ && \
	cd server.gen/ && \
	cfssl gencert \
		-ca=./../ca.gen/ca.pem \
		-ca-key=../ca.gen/ca-key.pem \
		-config=../ca-config.json \
		-profile=server \
		./../server.json | cfssljson -bare server && \
	cd ../
> @echo "Generating client keys and certificates"
> @mkdir -p client.gen/ && \
	cd client.gen/ && \
	cfssl gencert \
		-ca=../ca.gen/ca.pem \
		-ca-key=../ca.gen/ca-key.pem  \
		-config=../ca-config.json \
		-profile=client \
		./../client.json | cfssljson -bare client && \
	cd ../

.PHONY: certs
