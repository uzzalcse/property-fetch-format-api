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
        "/properties": {
            "get": {
                "description": "Get a list of properties based on provided property IDs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Comma-separated list of property IDs",
                        "name": "propertyIds",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The property list",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.PropertyResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
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
        "/v1/api/property/details/{propertyId}": {
            "get": {
                "description": "Get details of a property by its ID and language code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Property ID in format 'XX-1234'",
                        "name": "propertyId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Language code (default: en)",
                        "name": "languageCode",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Property details response",
                        "schema": {
                            "$ref": "#/definitions/models.PropertyResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request parameters",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Property not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
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
        "/v1/api/property/{propertyId}/gallery": {
            "get": {
                "description": "Retrieve property gallery images grouped by labels.",
                "tags": [
                    "Property Gallery"
                ],
                "summary": "Get property gallery images",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Property ID in format XX-123",
                        "name": "propertyId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "en",
                        "description": "Language code for the images",
                        "name": "languageCode",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Grouped images by label",
                        "schema": {
                            "$ref": "#/definitions/controllers.GroupedImages"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
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
        "/v1/api/user/": {
            "post": {
                "description": "Create a new user",
                "parameters": [
                    {
                        "description": "User details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.CreateUser"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "409": {
                        "description": "Email already exists",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/v1/api/user/{identifier}": {
            "get": {
                "description": "Retrieve user by ID or email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID (integer) or email (string)",
                        "name": "identifier",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "put": {
                "description": "Update user details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID or email",
                        "name": "identifier",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated user details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Validation error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete user by ID or email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID or email",
                        "name": "identifier",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.GroupedImages": {
            "type": "object",
            "additionalProperties": {
                "type": "array",
                "items": {
                    "type": "string"
                }
            }
        },
        "models.Category": {
            "type": "object",
            "properties": {
                "Display": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "LocationID": {
                    "type": "string"
                },
                "Name": {
                    "type": "string"
                },
                "Slug": {
                    "type": "string"
                },
                "Type": {
                    "type": "string"
                }
            }
        },
        "models.Counts": {
            "type": "object",
            "properties": {
                "Bathroom": {
                    "type": "integer"
                },
                "Bedroom": {
                    "type": "integer"
                },
                "Occupancy": {
                    "type": "integer"
                },
                "Reviews": {
                    "type": "integer"
                }
            }
        },
        "models.CreateUser": {
            "type": "object",
            "properties": {
                "Age": {
                    "type": "integer"
                },
                "Email": {
                    "type": "string"
                },
                "Name": {
                    "type": "string"
                }
            }
        },
        "models.GeoInfo": {
            "type": "object",
            "properties": {
                "Categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Category"
                    }
                },
                "City": {
                    "type": "string"
                },
                "Country": {
                    "type": "string"
                },
                "CountryCode": {
                    "type": "string"
                },
                "Display": {
                    "type": "string"
                },
                "Lat": {
                    "type": "string"
                },
                "Lng": {
                    "type": "string"
                },
                "LocationID": {
                    "type": "string"
                },
                "StateAbbr": {
                    "type": "string"
                }
            }
        },
        "models.Image": {
            "type": "object",
            "properties": {
                "Count": {
                    "type": "integer"
                },
                "Images": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.Partner": {
            "type": "object",
            "properties": {
                "Archived": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "BrandId": {
                    "type": "string"
                },
                "EpCluster": {
                    "type": "string"
                },
                "HcomID": {
                    "type": "string"
                },
                "ID": {
                    "type": "string"
                },
                "OwnerID": {
                    "type": "string"
                },
                "URL": {
                    "type": "string"
                },
                "UnitNumber": {
                    "type": "string"
                }
            }
        },
        "models.Property": {
            "type": "object",
            "properties": {
                "Amenities": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "Counts": {
                    "$ref": "#/definitions/models.Counts"
                },
                "EcoFriendly": {
                    "type": "boolean"
                },
                "FeatureImage": {
                    "type": "string"
                },
                "Image": {
                    "$ref": "#/definitions/models.Image"
                },
                "MinStay": {
                    "type": "integer"
                },
                "Price": {
                    "type": "number"
                },
                "PropertyName": {
                    "type": "string"
                },
                "PropertySlug": {
                    "type": "string"
                },
                "PropertyType": {
                    "type": "string"
                },
                "PropertyTypeCategoryId": {
                    "type": "string"
                },
                "ReviewScore": {
                    "type": "number"
                },
                "RoomSize": {
                    "type": "integer"
                },
                "UpdatedAt": {
                    "type": "string"
                }
            }
        },
        "models.PropertyResponse": {
            "type": "object",
            "properties": {
                "Feed": {
                    "type": "integer"
                },
                "GeoInfo": {
                    "$ref": "#/definitions/models.GeoInfo"
                },
                "ID": {
                    "type": "string"
                },
                "Partner": {
                    "$ref": "#/definitions/models.Partner"
                },
                "Property": {
                    "$ref": "#/definitions/models.Property"
                },
                "Published": {
                    "type": "boolean"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
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
