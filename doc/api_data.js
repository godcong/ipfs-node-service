define({ "api": [
  {
    "type": "get",
    "url": "/v1/info/:id",
    "title": "获取视频转换状态",
    "name": "info",
    "group": "Info",
    "version": "0.0.1",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "id",
            "description": "<p>文件名ID</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "http://localhost:8080/v1/status/9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm",
          "type": "string"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "string",
            "optional": false,
            "field": "code",
            "description": "<p>返回状态码：【正常：0】，【处理中：1】，【ID不存在：2】</p>"
          },
          {
            "group": "Success 200",
            "type": "string",
            "optional": false,
            "field": "msg",
            "description": "<p>返回具体消息</p>"
          },
          {
            "group": "Success 200",
            "type": "json",
            "optional": true,
            "field": "detail",
            "description": "<p>如正常则返回detail</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response OK:",
          "content": "{\n      \"code\":0,\n      \"msg\":\"ok\",\n}",
          "type": "json"
        },
        {
          "title": "Success-Response Processing:",
          "content": "{\n      \"code\":1,\n      \"msg\":\"processing\",\n}",
          "type": "json"
        }
      ]
    },
    "sampleRequest": [
      {
        "url": "/v1/status/:id"
      }
    ],
    "filename": "service/controller.go",
    "groupTitle": "Info",
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": "{\n  \"code\":-1,\n  \"msg\":\"error message\",\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "get",
    "url": "/v1/list",
    "title": "获取所有视频列表",
    "name": "info",
    "group": "Info",
    "version": "0.0.1",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "string",
            "optional": false,
            "field": "code",
            "description": "<p>返回状态码：【正常：0】，【处理中：1】，【ID不存在：2】</p>"
          },
          {
            "group": "Success 200",
            "type": "string",
            "optional": false,
            "field": "msg",
            "description": "<p>返回具体消息</p>"
          },
          {
            "group": "Success 200",
            "type": "json",
            "optional": true,
            "field": "detail",
            "description": "<p>如正常则返回detail</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response OK:",
          "content": "{\n      \"code\":0,\n      \"msg\":\"ok\",\n}",
          "type": "json"
        },
        {
          "title": "Success-Response Processing:",
          "content": "{\n      \"code\":1,\n      \"msg\":\"processing\",\n}",
          "type": "json"
        }
      ]
    },
    "sampleRequest": [
      {
        "url": "/v1/status/:id"
      }
    ],
    "parameter": {
      "examples": [
        {
          "title": "Request-Example:",
          "content": "http://localhost:8080/v1/list",
          "type": "string"
        }
      ]
    },
    "filename": "service/controller.go",
    "groupTitle": "Info",
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": "{\n  \"code\":-1,\n  \"msg\":\"error message\",\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/transfer",
    "title": "视频转换处理",
    "name": "transfer",
    "group": "Transfer",
    "version": "0.0.1",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "Binary",
            "optional": false,
            "field": "binary",
            "description": "<p>媒体文件二进制文件</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "{\n    \"id\":\"9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm\"\n}",
          "type": "string"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "    {\n      \"code\":0,\n      \"msg\":\"ok\",\n      \"detail\":{\n\t\t\t\"id\":\"9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm\"\n\t\t }\n    }",
          "type": "json"
        }
      ],
      "fields": {
        "detail": [
          {
            "group": "detail",
            "type": "string",
            "optional": false,
            "field": "id",
            "description": "<p>文件名ID</p>"
          }
        ],
        "Success 200": [
          {
            "group": "Success 200",
            "type": "string",
            "optional": false,
            "field": "msg",
            "description": "<p>返回具体消息</p>"
          },
          {
            "group": "Success 200",
            "type": "int",
            "optional": false,
            "field": "code",
            "description": "<p>返回状态码：【正常：0】，【失败，-1】</p>"
          },
          {
            "group": "Success 200",
            "type": "json",
            "optional": true,
            "field": "detail",
            "description": "<p>如正常则返回detail</p>"
          }
        ]
      }
    },
    "filename": "service/controller.go",
    "groupTitle": "Transfer",
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": "{\n  \"code\":-1,\n  \"msg\":\"error message\",\n}",
          "type": "json"
        }
      ]
    }
  },
  {
    "type": "post",
    "url": "/v1/upload",
    "title": "文件上传接口",
    "name": "upload",
    "group": "Upload",
    "version": "0.0.1",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "Binary",
            "optional": false,
            "field": "binary",
            "description": "<p>媒体文件二进制文件</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "upload a binary file from local",
          "type": "Binary"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "    {\n      \"code\":0,\n      \"msg\":\"ok\",\n      \"detail\":{\n\t\t\t\"id\":\"9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm\"\n\t\t }\n    }",
          "type": "json"
        }
      ],
      "fields": {
        "detail": [
          {
            "group": "detail",
            "type": "string",
            "optional": false,
            "field": "id",
            "description": "<p>文件名ID</p>"
          }
        ],
        "Success 200": [
          {
            "group": "Success 200",
            "type": "string",
            "optional": false,
            "field": "msg",
            "description": "<p>返回具体消息</p>"
          },
          {
            "group": "Success 200",
            "type": "int",
            "optional": false,
            "field": "code",
            "description": "<p>返回状态码：【正常：0】，【失败，-1】</p>"
          },
          {
            "group": "Success 200",
            "type": "json",
            "optional": true,
            "field": "detail",
            "description": "<p>如正常则返回detail</p>"
          }
        ]
      }
    },
    "sampleRequest": [
      {
        "url": "/v1/upload"
      }
    ],
    "filename": "service/controller.go",
    "groupTitle": "Upload",
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": "{\n  \"code\":-1,\n  \"msg\":\"error message\",\n}",
          "type": "json"
        }
      ]
    }
  }
] });
