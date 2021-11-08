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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/download/{id}": {
            "get": {
                "description": "download file from service",
                "produces": [
                    "*/*",
                    "image/jpeg",
                    "image/png",
                    "image/gif",
                    "video/*",
                    "audio/*",
                    "image/*",
                    "application/pdf",
                    "application/msword",
                    "application/vnd.ms-excel"
                ],
                "tags": [
                    "files"
                ],
                "summary": "download file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of the file",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "if file id is not valid",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "if file not found after upload",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/files": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return files info\nif you are not admin you can get info only about you files",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "get files",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of the user which files you want get",
                        "name": "user",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "sort by ascendig; can be name or date ",
                        "name": "sorted_by",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/file.File"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/files/upload": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "upload file to service",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "upload file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "file that need to upload",
                        "name": "uploadingForm",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "file description",
                        "name": "fileDescription",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/files.UploadFileResp"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "if file not found after upload",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/files/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get info about file",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "get file info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of the file",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/files.GetFileResp"
                        }
                    },
                    "400": {
                        "description": "if file id is not valid",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "if file not found after upload",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "delete file from service\nif you not admin you can only delete files that you upload",
                "tags": [
                    "files"
                ],
                "summary": "delete file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of the file",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "if file id is not valid",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "403": {
                        "description": "if it's not your file",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "if file not found after upload",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "err.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "file.File": {
            "type": "object",
            "properties": {
                "chunkSize": {
                    "type": "integer"
                },
                "filename": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "length": {
                    "type": "integer"
                },
                "metadata": {
                    "$ref": "#/definitions/file.Metadata"
                },
                "uploadDate": {
                    "type": "string"
                }
            }
        },
        "file.Metadata": {
            "type": "object",
            "properties": {
                "fileDescription": {
                    "type": "string"
                },
                "fileSender": {
                    "type": "string"
                }
            }
        },
        "files.GetFileResp": {
            "type": "object",
            "properties": {
                "chunkSize": {
                    "type": "integer"
                },
                "filename": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "length": {
                    "type": "integer"
                },
                "metadata": {
                    "$ref": "#/definitions/file.Metadata"
                },
                "uploadDate": {
                    "type": "string"
                }
            }
        },
        "files.UploadFileResp": {
            "type": "object",
            "properties": {
                "chunkSize": {
                    "type": "integer"
                },
                "filename": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "length": {
                    "type": "integer"
                },
                "metadata": {
                    "$ref": "#/definitions/file.Metadata"
                },
                "uploadDate": {
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
	Version:     "1.0",
	Host:        "",
	BasePath:    "/api/mfs",
	Schemes:     []string{},
	Title:       "MicroFileService API",
	Description: "This is a server for save and get files",
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
