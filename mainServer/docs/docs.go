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
        "/articles/{articleID}/versions/{versionID}": {
            "post": {
                "description": "Upload files to update an article version, can only be done by an owner. Requires multipart form data, with a file attached as the field \"file\"",
                "consumes": [
                    "multipart/form-data"
                ],
                "summary": "Update article version",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Article ID",
                        "name": "articleID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Version ID",
                        "name": "versionID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success"
                    },
                    "400": {
                        "description": "Bad request, possibly bad file data or permissions"
                    },
                    "404": {
                        "description": "Specified article version not found"
                    }
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
