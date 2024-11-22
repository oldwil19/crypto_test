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
        "/account/balance/add": {
            "post": {
                "description": "Permite añadir saldo en USD al balance del usuario autenticado. Requiere un token JWT válido.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Añadir saldo al usuario",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003ctoken\u003e",
                        "description": "Token JWT",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Cuerpo de la solicitud para añadir saldo, con el monto en USD",
                        "name": "AddBalanceRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/application.AddBalanceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Respuesta exitosa con mensaje y balance actualizado",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Error de validación: La cantidad ingresada es inválida",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Error de usuario: Usuario no encontrado",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Error interno: No se pudo actualizar la información del usuario",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Permite a los usuarios autenticarse y obtener un token JWT para acceder a los endpoints protegidos.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Autenticación de usuarios",
                "parameters": [
                    {
                        "description": "Credenciales del usuario para autenticación",
                        "name": "LoginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/application.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Token JWT generado para el usuario autenticado",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Error en la validación de datos enviados",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Credenciales inválidas o usuario no encontrado",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Error interno al generar el token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "application.AddBalanceRequest": {
            "type": "object",
            "required": [
                "amount"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                }
            }
        },
        "application.LoginRequest": {
            "description": "Estructura del cuerpo de la solicitud para el endpoint de login.",
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "description": "La contraseña debe ser obligatoria.",
                    "type": "string"
                },
                "username": {
                    "description": "El nombre de usuario debe ser obligatorio.",
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
