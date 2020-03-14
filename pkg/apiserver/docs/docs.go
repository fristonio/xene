// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Deepesh Pathak",
            "url": "https://dpathak.co",
            "email": "deepeshpathak09@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/authok": {
            "get": {
                "description": "auth OK handler handles the ping to api routes which are",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api"
                ],
                "summary": "Handles ping event for api routes.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPMessage"
                        }
                    }
                }
            }
        },
        "/api/v1/registry/workflow": {
            "get": {
                "description": "If a name is provided return the corresponding workflow object, if prefix  is set to some value",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "registry"
                ],
                "summary": "Returns the specified workflow object from the store.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Prefix based get for workflow.",
                        "name": "prefix",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "name of the workflow to get.",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Workflow"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "description": "This route creates a new workflow for xene to operate on, if the workflow already exists",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "registry"
                ],
                "summary": "Creates a new workflow in the store.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Workflow manifest to be created.",
                        "name": "workflow",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/registry/workflow/{name}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "registry"
                ],
                "summary": "Returns the specified workflow object from the store with the name in params.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of the workflow to get.",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Workflow"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes the workflow specified by the name parameter, if the workflow is not",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "registry"
                ],
                "summary": "Deletes the specified workflow from the store.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the workflow to be deleted.",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/:provider": {
            "get": {
                "description": "Log in to xene using the configured oauth providers that xene supports.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Handles login for xene",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider for oauth login",
                        "name": "provider",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.OauthLogin"
                        }
                    }
                }
            }
        },
        "/auth/:provider/redirect": {
            "get": {
                "description": "redirectHandler handles the redirect from the Oauth provider after the authentication process has",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Handles redirect from the login oauth provider.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider for the oauth redirect",
                        "name": "provider",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JWTAuth"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/refresh/": {
            "get": {
                "description": "Handles authentication token refresh",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Handle authentication token refresh for the oauth provider.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JWTAuth"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response.HTTPError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Invalid authentication type provided."
                }
            }
        },
        "response.HTTPMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Messsage in response to your request"
                }
            }
        },
        "response.JWTAuth": {
            "type": "object",
            "properties": {
                "expiresIn": {
                    "type": "string",
                    "example": "3600"
                },
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
                },
                "userEmail": {
                    "type": "string",
                    "example": "example@example.com"
                },
                "userName": {
                    "type": "string",
                    "example": "fristonio"
                }
            }
        },
        "response.OauthLogin": {
            "type": "object",
            "properties": {
                "loginURL": {
                    "description": "LoginURL is the URL to be used for logging in.",
                    "type": "string",
                    "example": "https://xxxx.io/login"
                }
            }
        },
        "response.Workflow": {
            "type": "object",
            "properties": {
                "workflow": {
                    "type": "string",
                    "example": "Workflow Document"
                }
            }
        },
        "response.WorkflowsFromPrefix": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer",
                    "example": 2
                },
                "workflows": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1.0",
	Host:        "localhost:6060",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Xene API server",
	Description: "Xene is the workflow creator and manager tool",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
