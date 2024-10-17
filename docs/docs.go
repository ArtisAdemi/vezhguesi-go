// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/api/auth/": {
            "post": {
                "description": "Validates email, username, first name, last name, password checks if email exists, if not creates new user and sends email with verification link.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Signup",
                "parameters": [
                    {
                        "description": "SignupRequest",
                        "name": "SignupRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.SignupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.SignupResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/login": {
            "post": {
                "description": "Validates email and password in request, check if user exists in DB if not throw 404 otherwise compare the request password with hash, then check if user is active, then finds relationships of user with orgs and then generates a JWT token, and returns UserData, Orgs, and Token in response.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "LoginRequest",
                        "name": "LoginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/verify-signup/{token}": {
            "get": {
                "description": "Validates token in param, if token parses valid then user will be verified and be updated in DB.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "VerifySignup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.StatusResponse"
                        }
                    }
                }
            }
        },
        "/api/reports/": {
            "get": {
                "description": "Validates user id. Gets all reports",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Get Reports",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key (e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/reports.ReportsResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Validates subject, start date, end date. Creates a new report.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Create Report",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key (e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "CreateReportRequest",
                        "name": "CreateReportRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/reports.CreateReportRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/reports.ReportResponse"
                        }
                    }
                }
            }
        },
        "/api/reports/{id}": {
            "get": {
                "description": "Validates id and user id. Gets report by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Get Report By ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key (e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Report ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/reports.ReportsResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Validates id and user id. Updates report",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Update Report",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key (e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Report ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "UpdateReportRequest",
                        "name": "UpdateReportRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/reports.UpdateReportRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/reports.ReportResponse"
                        }
                    }
                }
            }
        },
        "/api/users": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "GetUsers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/users.UserResponse"
                            }
                        }
                    }
                }
            }
        },
        "/api/users/{userId}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "GetUserByID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key (e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.FindByIDResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "userData": {
                    "$ref": "#/definitions/auth.UserData"
                }
            }
        },
        "auth.SignupRequest": {
            "type": "object",
            "properties": {
                "confirmPassword": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "auth.SignupResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "auth.StatusResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
                }
            }
        },
        "auth.UserData": {
            "type": "object",
            "properties": {
                "avatarImgUrl": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "profileId": {
                    "type": "integer"
                },
                "role": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "reports.CreateReportRequest": {
            "type": "object",
            "properties": {
                "endDate": {
                    "type": "string"
                },
                "startDate": {
                    "type": "string"
                },
                "subject": {
                    "type": "string"
                }
            }
        },
        "reports.Report": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "endDate": {
                    "type": "string"
                },
                "entities": {
                    "type": "string"
                },
                "findings": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "reportText": {
                    "type": "string"
                },
                "sentiment": {
                    "type": "integer"
                },
                "sourceID": {
                    "type": "integer"
                },
                "startDate": {
                    "type": "string"
                },
                "subject": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "reports.ReportResponse": {
            "type": "object",
            "properties": {
                "report": {
                    "$ref": "#/definitions/reports.Report"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "reports.ReportsResponse": {
            "type": "object",
            "properties": {
                "reports": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/reports.Report"
                    }
                }
            }
        },
        "reports.UpdateReportRequest": {
            "type": "object",
            "properties": {
                "entities": {
                    "type": "string"
                },
                "findings": {
                    "type": "string"
                },
                "reportText": {
                    "type": "string"
                },
                "sentiment": {
                    "type": "integer"
                },
                "sourceId": {
                    "type": "integer"
                },
                "subject": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "users.FindByIDResponse": {
            "type": "object",
            "properties": {
                "avatarImgUrl": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "users.UserResponse": {
            "type": "object",
            "properties": {
                "avatarImgUrl": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
