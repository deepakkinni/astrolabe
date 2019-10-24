// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Arachne API",
    "version": "1.0.0"
  },
  "paths": {
    "/arachne": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "List available services",
        "operationId": "listServices",
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ServiceList"
            }
          }
        }
      }
    },
    "/arachne/tasks": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Lists running and recent tasks",
        "operationId": "listTasks",
        "responses": {
          "200": {
            "description": "List of recent task IDs",
            "schema": {
              "$ref": "#/definitions/TaskIDList"
            }
          }
        }
      }
    },
    "/arachne/tasks/{taskID}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Gets info about a running or recently completed task",
        "operationId": "getTaskInfo",
        "parameters": [
          {
            "type": "string",
            "description": "The ID of the task to retrieve info for",
            "name": "taskID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Info for running or recently completed task",
            "schema": {
              "$ref": "#/definitions/TaskInfo"
            }
          }
        }
      }
    },
    "/arachne/{service}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "List protected entities for the service.  Results will be returned in canonical ID order (string sorted).  Fewer results may be returned than expected, the ProtectedEntityList has a field specifying if the list has been truncated.",
        "operationId": "listProtectedEntities",
        "parameters": [
          {
            "type": "string",
            "description": "The service to list protected entities from",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "integer",
            "format": "int32",
            "description": "The maximum number of results to return (fewer results may be returned)",
            "name": "maxResults",
            "in": "query"
          },
          {
            "type": "string",
            "description": "Results will be returned that come after this ID",
            "name": "idsAfter",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityList"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Copy a protected entity into the repository.  There is no option to embed data on this path, for a self-contained or partially self-contained object, use the restore from zip file option in the S3 API REST API",
        "operationId": "copyProtectedEntity",
        "parameters": [
          {
            "type": "string",
            "description": "The service to copy the protected entity into",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "create",
              "create_new",
              "update"
            ],
            "type": "string",
            "description": "How to handle the copy.  create - a new protected entity with the Protected Entity ID will be created.  If the Protected Entity ID already exists, the copy will fail.  create_new - A Protected Entity with a new ID will be created with data and metadata from the source protected entity.  Update - If a protected entity with the same ID exists it will be overwritten.  If there is no PE with that ID, one will be created with the same ID. For complex Persistent Entities, the mode will be applied to all of the component entities that are part of this operation as well.",
            "name": "mode",
            "in": "query",
            "required": true
          },
          {
            "description": "Info of ProtectedEntity to copy",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/ProtectedEntityInfo"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created - returned if the protected entity can be created immediately",
            "schema": {
              "$ref": "#/definitions/CreatedResponse"
            }
          },
          "202": {
            "description": "Create in progress",
            "schema": {
              "$ref": "#/definitions/CreateInProgressResponse"
            }
          }
        }
      }
    },
    "/arachne/{service}/{protectedEntityID}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Get the info for a Protected Entity including name, data access and components",
        "operationId": "getProtectedEntityInfo",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to retrieve info for",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityInfo"
            }
          }
        }
      },
      "delete": {
        "produces": [
          "application/json"
        ],
        "summary": "Deletes a protected entity or snapshot of a protected entity (if the snapshot ID is specified)",
        "operationId": "deleteProtectedEntity",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to retrieve info for",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityID"
            }
          }
        }
      }
    },
    "/arachne/{service}/{protectedEntityID}/snapshots": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Gets the list of snapshots for this protected entity",
        "operationId": "getSnapshots",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to retrieve info for",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "\"200\" response\"",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityList"
            }
          }
        }
      },
      "post": {
        "produces": [
          "application/json"
        ],
        "summary": "Creates a new snapshot for this protected entity",
        "operationId": "createSnapshot",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to snapshot",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Snapshot created successfully, returns the new snapshot ID",
            "schema": {
              "$ref": "#/definitions/ProtectedEntitySnapshotID"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "ComponentSpec": {
      "type": "object",
      "required": [
        "id"
      ],
      "properties": {
        "id": {
          "$ref": "#/definitions/ProtectedEntityID"
        },
        "server": {
          "type": "string"
        }
      }
    },
    "CreateInProgressResponse": {
      "type": "object",
      "properties": {
        "taskID": {
          "type": "string"
        }
      }
    },
    "CreatedResponse": {
      "type": "string"
    },
    "DataTransport": {
      "type": "object",
      "properties": {
        "params": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "transportType": {
          "type": "string"
        }
      }
    },
    "ProtectedEntityID": {
      "type": "string"
    },
    "ProtectedEntityInfo": {
      "type": "object",
      "required": [
        "id",
        "name",
        "metadataTransports",
        "dataTransports",
        "combinedTransports",
        "componentSpecs"
      ],
      "properties": {
        "combinedTransports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/DataTransport"
          }
        },
        "componentSpecs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ComponentSpec"
          }
        },
        "dataTransports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/DataTransport"
          }
        },
        "id": {
          "$ref": "#/definitions/ProtectedEntityID"
        },
        "metadataTransports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/DataTransport"
          }
        },
        "name": {
          "type": "string"
        }
      }
    },
    "ProtectedEntityList": {
      "type": "object",
      "properties": {
        "list": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ProtectedEntityID"
          }
        },
        "truncated": {
          "type": "boolean"
        }
      }
    },
    "ProtectedEntitySnapshotID": {
      "type": "string"
    },
    "ServiceList": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "TaskID": {
      "type": "string"
    },
    "TaskIDList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/TaskID"
      }
    },
    "TaskInfo": {
      "type": "object",
      "required": [
        "id",
        "completed",
        "status",
        "startedTime",
        "progress"
      ],
      "properties": {
        "completed": {
          "type": "boolean"
        },
        "details": {
          "type": "string"
        },
        "finishedTime": {
          "type": "string"
        },
        "id": {
          "$ref": "#/definitions/TaskID"
        },
        "progress": {
          "type": "integer"
        },
        "startedTime": {
          "type": "string"
        },
        "status": {
          "type": "string",
          "enum": [
            "running",
            "success",
            "failed",
            "cancelled"
          ]
        }
      }
    }
  },
  "x-components": {}
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Arachne API",
    "version": "1.0.0"
  },
  "paths": {
    "/arachne": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "List available services",
        "operationId": "listServices",
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ServiceList"
            }
          }
        }
      }
    },
    "/arachne/tasks": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Lists running and recent tasks",
        "operationId": "listTasks",
        "responses": {
          "200": {
            "description": "List of recent task IDs",
            "schema": {
              "$ref": "#/definitions/TaskIDList"
            }
          }
        }
      }
    },
    "/arachne/tasks/{taskID}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Gets info about a running or recently completed task",
        "operationId": "getTaskInfo",
        "parameters": [
          {
            "type": "string",
            "description": "The ID of the task to retrieve info for",
            "name": "taskID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Info for running or recently completed task",
            "schema": {
              "$ref": "#/definitions/TaskInfo"
            }
          }
        }
      }
    },
    "/arachne/{service}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "List protected entities for the service.  Results will be returned in canonical ID order (string sorted).  Fewer results may be returned than expected, the ProtectedEntityList has a field specifying if the list has been truncated.",
        "operationId": "listProtectedEntities",
        "parameters": [
          {
            "type": "string",
            "description": "The service to list protected entities from",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "integer",
            "format": "int32",
            "description": "The maximum number of results to return (fewer results may be returned)",
            "name": "maxResults",
            "in": "query"
          },
          {
            "type": "string",
            "description": "Results will be returned that come after this ID",
            "name": "idsAfter",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityList"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Copy a protected entity into the repository.  There is no option to embed data on this path, for a self-contained or partially self-contained object, use the restore from zip file option in the S3 API REST API",
        "operationId": "copyProtectedEntity",
        "parameters": [
          {
            "type": "string",
            "description": "The service to copy the protected entity into",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "create",
              "create_new",
              "update"
            ],
            "type": "string",
            "description": "How to handle the copy.  create - a new protected entity with the Protected Entity ID will be created.  If the Protected Entity ID already exists, the copy will fail.  create_new - A Protected Entity with a new ID will be created with data and metadata from the source protected entity.  Update - If a protected entity with the same ID exists it will be overwritten.  If there is no PE with that ID, one will be created with the same ID. For complex Persistent Entities, the mode will be applied to all of the component entities that are part of this operation as well.",
            "name": "mode",
            "in": "query",
            "required": true
          },
          {
            "description": "Info of ProtectedEntity to copy",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/ProtectedEntityInfo"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created - returned if the protected entity can be created immediately",
            "schema": {
              "$ref": "#/definitions/CreatedResponse"
            }
          },
          "202": {
            "description": "Create in progress",
            "schema": {
              "$ref": "#/definitions/CreateInProgressResponse"
            }
          }
        }
      }
    },
    "/arachne/{service}/{protectedEntityID}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Get the info for a Protected Entity including name, data access and components",
        "operationId": "getProtectedEntityInfo",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to retrieve info for",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityInfo"
            }
          }
        }
      },
      "delete": {
        "produces": [
          "application/json"
        ],
        "summary": "Deletes a protected entity or snapshot of a protected entity (if the snapshot ID is specified)",
        "operationId": "deleteProtectedEntity",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to retrieve info for",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "200 response",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityID"
            }
          }
        }
      }
    },
    "/arachne/{service}/{protectedEntityID}/snapshots": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Gets the list of snapshots for this protected entity",
        "operationId": "getSnapshots",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to retrieve info for",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "\"200\" response\"",
            "schema": {
              "$ref": "#/definitions/ProtectedEntityList"
            }
          }
        }
      },
      "post": {
        "produces": [
          "application/json"
        ],
        "summary": "Creates a new snapshot for this protected entity",
        "operationId": "createSnapshot",
        "parameters": [
          {
            "type": "string",
            "description": "The service for the protected entity",
            "name": "service",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "The protected entity ID to snapshot",
            "name": "protectedEntityID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Snapshot created successfully, returns the new snapshot ID",
            "schema": {
              "$ref": "#/definitions/ProtectedEntitySnapshotID"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "ComponentSpec": {
      "type": "object",
      "required": [
        "id"
      ],
      "properties": {
        "id": {
          "$ref": "#/definitions/ProtectedEntityID"
        },
        "server": {
          "type": "string"
        }
      }
    },
    "CreateInProgressResponse": {
      "type": "object",
      "properties": {
        "taskID": {
          "type": "string"
        }
      }
    },
    "CreatedResponse": {
      "type": "string"
    },
    "DataTransport": {
      "type": "object",
      "properties": {
        "params": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "transportType": {
          "type": "string"
        }
      }
    },
    "ProtectedEntityID": {
      "type": "string"
    },
    "ProtectedEntityInfo": {
      "type": "object",
      "required": [
        "id",
        "name",
        "metadataTransports",
        "dataTransports",
        "combinedTransports",
        "componentSpecs"
      ],
      "properties": {
        "combinedTransports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/DataTransport"
          }
        },
        "componentSpecs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ComponentSpec"
          }
        },
        "dataTransports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/DataTransport"
          }
        },
        "id": {
          "$ref": "#/definitions/ProtectedEntityID"
        },
        "metadataTransports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/DataTransport"
          }
        },
        "name": {
          "type": "string"
        }
      }
    },
    "ProtectedEntityList": {
      "type": "object",
      "properties": {
        "list": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ProtectedEntityID"
          }
        },
        "truncated": {
          "type": "boolean"
        }
      }
    },
    "ProtectedEntitySnapshotID": {
      "type": "string"
    },
    "ServiceList": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "TaskID": {
      "type": "string"
    },
    "TaskIDList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/TaskID"
      }
    },
    "TaskInfo": {
      "type": "object",
      "required": [
        "id",
        "completed",
        "status",
        "startedTime",
        "progress"
      ],
      "properties": {
        "completed": {
          "type": "boolean"
        },
        "details": {
          "type": "string"
        },
        "finishedTime": {
          "type": "string"
        },
        "id": {
          "$ref": "#/definitions/TaskID"
        },
        "progress": {
          "type": "integer"
        },
        "startedTime": {
          "type": "string"
        },
        "status": {
          "type": "string",
          "enum": [
            "running",
            "success",
            "failed",
            "cancelled"
          ]
        }
      }
    }
  },
  "x-components": {}
}`))
}