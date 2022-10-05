# Go Modular REST API (POC w/ DynamoDB)

- [Go Modular REST API (POC w/ DynamoDB)](#go-modular-rest-api-poc-w-dynamodb)
  - [Testing](#testing)
  - [How it works](#how-it-works)
    - [Entities](#entities)
    - [Rules](#rules)
    - [Adapters](#adapters)
    - [Routes, Handlers, and Interfaces](#routes-handlers-and-interfaces)
    - [Healthcheck, HTTP responses, and Logging](#healthcheck-http-responses-and-logging)
    - [A few noteworthy functions](#a-few-noteworthy-functions)
    - [CORS](#cors)

---

> **NOTE: This repo is a POC.**

While the structure seems monolithic, the objective is to build a modular REST API out of submodules.

This API will easily scale to 10k concurrent connections with zero issues.

The modularity allows a concurrent separation of concerns with multiple levels of abstraction, particularly due to a heavy use of interfaces.

---

## Testing

1. Make sure you have `~/.aws/credentials` and `~/.aws/config` configured.

2. At the project root, run `go run ./app/main.go`.

3. Use curl, Postman, or whatever tool you prefer to test routes.

---

## How it works

This API publishes at port 8080; configure in [config.go](/config/config.go).

When the API starts, it opens a session to DynamoDB by loading credentials from `~/.aws/credentials` and the region from `~/.aws/config`.

The following routes are configured:

- `/healthcheck`
- `/object`

### Entities

**Entities** are a represenation of DynamoDB objects and have two parts:

- **Base** which is the baseline metadata for an object such as its ID, creation and update timestamp.
- The actual **Object**, which consists of its Base and its Name.

### Rules

**Rules** perform operations such as marshalling/unmarshalling returned data to a struct, creating new tables, and validating data before its processed. When data is marshalled/unmarshalled, it is processed to/from **JSON** via structs.

### Adapters

Adapters handle database table functions.

### Routes, Handlers, and Interfaces

When a route is accessed, its Handler hands off operations to Interfaces, and the Interfaces perform the database operations.

The following Interfaces are configured in to parts: Control level interfaces and Adapter level interfaces. The Control level functions rely on the Adapter level functions in a one-to-one relationship.

Control interface functions:

- **Create()** — standard Create operation.
- **DescribeOne()** — standard Read operation; applies to a single object.
- **DescribeAll()** — standard Read operation; applies to all objects.
- **Update()** — standard Update operation.
- **Remove()** — standard Delete operation.

Adapter interface functions:

- **Create** — standard Create operation.
- **FindOne()** — standard Read operation; applies to a single object.
- **FindAll()** — standard Read operation; applies to all objects.
- **Update** — standard Update operation.
- **Delete** — standard Delete operation.

### Healthcheck, HTTP responses, and Logging

Healthcheck is configured at [healthcheck.go](/internal/handlers/healthcheck/healthcheck.go) with the following HTTP statuses:

| HTTP Status | Description |
| ----------- | ----------- |
| 200 | Service OK |
| 204 | No Content |
| 400 | Bad Request |
| 404 | Not Found |
| 405 | Method Not Allowed |
| 409 | Conflict |
| 500 | Internal Server Error |

HTTP responses are configured at [response.go](/utils/http/response.go).

Logging is configured at [logger.go](/utils/logger/logger.go).

### A few noteworthy functions

- **Migrate()** in [object.go](/internal/rules/object/object.go) allows you to migrate your tables.
- **checkTables()** in [main.go](/app/main.go) performs error handling to check that the DynamoDB table exists.

### CORS

**Cross-Origin Resource Sharing (CORS)** is enabled by default; configure in [routes.go](/internal/routes/routes.go).
