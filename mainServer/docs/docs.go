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
        "/articles": {
            "post": {
                "description": "Creates new article, including main article version. Returns main version. Owners must be specified as email addresses, not usernames.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create new article",
                "parameters": [
                    {
                        "description": "Article info",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ArticleCreationForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Version"
                        }
                    }
                }
            }
        },
        "/articles/{articleID}/versions": {
            "get": {
                "description": "Gets all versions belonging to a specific article. Does not include version contents.",
                "produces": [
                    "application/json"
                ],
                "summary": "List article versions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Article ID",
                        "name": "articleID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Version"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request input"
                    },
                    "500": {
                        "description": "Could not get versions from server"
                    }
                }
            }
        },
        "/articles/{articleID}/versions/{versionID}": {
            "get": {
                "description": "Gets the version content + metadata from the database + filesystem. Must be accessible without being authenticated.",
                "produces": [
                    "application/json"
                ],
                "summary": "Get version content + metadata",
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
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Version"
                        }
                    },
                    "404": {
                        "description": "Version not found"
                    }
                }
            },
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
        },
        "/createExampleUser": {
            "post": {
                "description": "Creates a hardcoded user entity and adds it to the database, demonstrates how to add to database",
                "produces": [
                    "text/plain"
                ],
                "summary": "Temporary user creation endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getExampleUser": {
            "get": {
                "description": "Returns a user with a hardcoded email address, demonstrates how to use the services.",
                "produces": [
                    "application/json"
                ],
                "summary": "Get test user from database endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.User"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "pwd": {
                    "type": "string"
                }
            }
        },
        "models.ArticleCreationForm": {
            "type": "object",
            "required": [
                "owners",
                "title"
            ],
            "properties": {
                "owners": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Version": {
            "type": "object",
            "properties": {
                "articleID": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "owners": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "versionID": {
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
