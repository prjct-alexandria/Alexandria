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
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    }
                }
            }
        },
        "/articles/{articleID}/requests": {
            "post": {
                "description": "Creates request to merge one article versions' changes into another",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create request",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Article ID",
                        "name": "articleID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RequestCreationForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Request"
                        }
                    },
                    "400": {
                        "description": "Invalid article ID or request creation data"
                    },
                    "500": {
                        "description": "Error creating request on server"
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
        },
        "/login": {
            "post": {
                "description": "Takes in user email and password from a JSON body, verifies if they are correct with the database and returns a JWT token",
                "consumes": [
                    "application/json"
                ],
                "summary": "Endpoint for user logging in",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success"
                    },
                    "400": {
                        "description": "Invalid JSON provided"
                    },
                    "403": {
                        "description": "Invalid password"
                    },
                    "500": {
                        "description": "Could not create token"
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Takes in user credentials from a JSON body, and makes sure they are securely stored in the database.",
                "consumes": [
                    "application/json"
                ],
                "summary": "Endpoint for user registration",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success"
                    },
                    "400": {
                        "description": "Invalid user JSON provided"
                    },
                    "403": {
                        "description": "Could not generate password hash"
                    },
                    "409": {
                        "description": "Could not save user to database"
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
        "httperror.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
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
        "models.LoginForm": {
            "type": "object",
            "required": [
                "email",
                "pwd"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "pwd": {
                    "type": "string"
                }
            }
        },
        "models.Request": {
            "type": "object",
            "properties": {
                "articleID": {
                    "type": "integer"
                },
                "requestID": {
                    "type": "integer"
                },
                "sourceHistoryID": {
                    "type": "string"
                },
                "sourceVersionID": {
                    "type": "integer"
                },
                "state": {
                    "type": "string"
                },
                "targetHistoryID": {
                    "type": "string"
                },
                "targetVersionID": {
                    "type": "integer"
                }
            }
        },
        "models.RequestCreationForm": {
            "type": "object",
            "required": [
                "sourceHistoryID",
                "sourceVersionID",
                "targetHistoryID",
                "targetVersionID"
            ],
            "properties": {
                "sourceHistoryID": {
                    "type": "string"
                },
                "sourceVersionID": {
                    "type": "integer"
                },
                "targetHistoryID": {
                    "type": "string"
                },
                "targetVersionID": {
                    "type": "integer"
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
