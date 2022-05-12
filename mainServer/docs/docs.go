// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/helloWorldJson": {
            "get": {
                "description": "Returns a list of joke in JSON format, demonstrates how future endpoints could work.",
                "produces": [
                    "application/json"
                ],
                "summary": "Hello world test endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.Joke"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Joke": {
            "type": "object",
            "required": [
                "id",
                "joke"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "joke": {
                    "type": "string"
                },
                "likes": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "API documentation",
	Description:      "Documentation of Alexandria's API. Endpoints can be tried out directly in this interactive documentation.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
