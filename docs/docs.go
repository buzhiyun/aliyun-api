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
        "/api/cdn/refresh": {
            "post": {
                "description": "根据提供的URL去刷新CDN",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cdn"
                ],
                "summary": "刷新CDN",
                "parameters": [
                    {
                        "description": "urls 是 []string",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.RefreshCdnReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    }
                }
            }
        },
        "/api/cms/ecs/cpu": {
            "post": {
                "description": "根据ECS  获取CPU使用率",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "监控CMS"
                ],
                "summary": "获取ECS CPU信息",
                "parameters": [
                    {
                        "description": "hostName 是 []string",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.getDataPointReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.ApiJson"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/cms.Datapoint"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    }
                }
            }
        },
        "/api/ecs/refresh": {
            "post": {
                "description": "刷新ecs实例列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ecs"
                ],
                "summary": "刷新ecs实例列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    }
                }
            }
        },
        "/api/ecs/search": {
            "post": {
                "description": "会设置服务器 所有负载均衡里的权重",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ecs"
                ],
                "summary": "设置所有包含该服务器的负载均衡里的权重",
                "parameters": [
                    {
                        "description": "hostname 必填 ；weight 为 0-100之间的数字 必填 ",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.setEcsWeight"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    }
                }
            }
        },
        "/api/slb/acl/add": {
            "post": {
                "description": "添加主机到ACL",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "slb"
                ],
                "summary": "添加主机到ACL",
                "parameters": [
                    {
                        "description": "acl_id 是 slb accessList的ID ； host 和 ip不能同时为空 ，类型均为 []string",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.AclListReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    }
                }
            }
        },
        "/api/slb/acl/delete": {
            "post": {
                "description": "从ACL移除主机",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "slb"
                ],
                "summary": "从ACL移除主机",
                "parameters": [
                    {
                        "description": "acl_id 是 slb accessList的ID ； host 和 ip不能同时为空 ，类型均为 []string",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.AclListReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiJson"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "cms.Datapoint": {
            "type": "object",
            "properties": {
                "Average": {
                    "type": "number"
                },
                "Maximum": {
                    "type": "number"
                },
                "Minimum": {
                    "type": "number"
                },
                "instanceId": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "controllers.AclListReq": {
            "type": "object",
            "required": [
                "acl_id"
            ],
            "properties": {
                "acl_id": {
                    "type": "string"
                },
                "comment": {
                    "type": "string"
                },
                "host": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "ip": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "controllers.RefreshCdnReq": {
            "type": "object",
            "required": [
                "urls"
            ],
            "properties": {
                "urls": {
                    "description": "主机名,支持通配符",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "controllers.getDataPointReq": {
            "type": "object",
            "properties": {
                "hostName": {
                    "description": "主机名,支持通配符",
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "sdcf_v3_030"
                    ]
                },
                "instanceId": {
                    "description": "主机Id",
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "i-bp1d1oh9a06r70buf03l"
                    ]
                }
            }
        },
        "controllers.searchHostReq": {
            "type": "object",
            "properties": {
                "hostname": {
                    "description": "主机名,支持通配符",
                    "type": "string"
                },
                "ip": {
                    "description": "主机名,支持通配符",
                    "type": "string"
                }
            }
        },
        "controllers.setEcsWeight": {
            "type": "object",
            "required": [
                "hostname",
                "weight"
            ],
            "properties": {
                "hostname": {
                    "description": "主机名,支持通配符",
                    "type": "string"
                },
                "weight": {
                    "description": "权重  *int 防止0在json的时候被丢掉",
                    "type": "integer",
                    "maximum": 100,
                    "minimum": 0
                }
            }
        },
        "utils.ApiJson": {
            "type": "object",
            "properties": {
                "data": {},
                "msg": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
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
