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
            "name": "xusheng",
            "url": "https://github.com/xuxusheng",
            "email": "20691718@qq.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/users": {
            "get": {
                "description": "通过 name、phone 字段查询匹配的用户，支持模糊查询、分页",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "分页获取用户列表",
                "parameters": [
                    {
                        "type": "string",
                        "default": "",
                        "description": "用户名",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "",
                        "description": "手机号",
                        "name": "phone",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "1",
                        "description": "第几页",
                        "name": "pn",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "10",
                        "description": "每页记录数量",
                        "name": "ps",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "allOf": [
                                                {
                                                    "$ref": "#/definitions/model.DWithP"
                                                },
                                                {
                                                    "type": "object",
                                                    "properties": {
                                                        "data": {
                                                            "type": "array",
                                                            "items": {
                                                                "$ref": "#/definitions/model.User"
                                                            }
                                                        }
                                                    }
                                                }
                                            ]
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "description": "创建新用户接口，专供管理平台调用",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "创建新用户",
                "parameters": [
                    {
                        "maxLength": 20,
                        "minLength": 6,
                        "description": "用户名（6-20位数字或字母构成)",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "maxLength": 11,
                        "minLength": 11,
                        "description": "手机号（十一位数字）",
                        "name": "phone",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "密码",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.User"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "内部错误",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResp"
                        }
                    }
                }
            }
        },
        "/api/v1/users/{id}": {
            "get": {
                "description": "通过 ID 查询单个用户详细信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "查询单个用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户ID",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.Resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.User"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "put": {
                "description": "更新用户名、手机号等，字段不传或为空字符串不修改此字段。\n只针对用户基本信息修改，其他信息例如角色、密码等，通过专门的接口改，便于权限控制。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "更新用户信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "用户名",
                        "name": "name",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "手机号",
                        "name": "phone",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Resp"
                        }
                    }
                }
            },
            "delete": {
                "description": "根据 ID 删除用户",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "删除用户",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Resp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.DWithP": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "pn": {
                    "description": "当前页码",
                    "type": "integer"
                },
                "ps": {
                    "description": "每页显示记录数",
                    "type": "integer"
                },
                "total": {
                    "description": "总共多少条记录",
                    "type": "integer"
                }
            }
        },
        "model.ErrMeta": {
            "type": "object",
            "properties": {
                "err_code": {
                    "description": "错误码",
                    "type": "integer",
                    "example": 10000000
                },
                "err_details": {
                    "description": "错误详情",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "err_msg": {
                    "description": "错误信息",
                    "type": "string",
                    "example": "服务器内部错误"
                }
            }
        },
        "model.ErrResp": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "meta": {
                    "$ref": "#/definitions/model.ErrMeta"
                }
            }
        },
        "model.Meta": {
            "type": "object",
            "properties": {
                "err_code": {
                    "description": "错误码",
                    "type": "integer",
                    "example": 0
                },
                "err_details": {
                    "description": "错误详情",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "err_msg": {
                    "description": "错误信息",
                    "type": "string",
                    "example": "成功"
                }
            }
        },
        "model.Resp": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "meta": {
                    "$ref": "#/definitions/model.Meta"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2020-12-09T18:52:41.555+08:00"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "xusheng"
                },
                "phone": {
                    "type": "string",
                    "example": "17707272442"
                },
                "role": {
                    "type": "string",
                    "example": "member"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2020-12-09T18:52:41.555+08:00"
                }
            }
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
	BasePath:    "",
	Schemes:     []string{},
	Title:       "时频学习平台",
	Description: "时频学习平台后端接口文档",
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
