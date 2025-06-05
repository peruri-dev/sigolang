# SiGolang Boilerplate

* **mux**: `fiber` over `humafiber`.
* **logging**: `inalog`.
* **cli**: `spf13/cobra` over `humacli`.
* **config**: `cleanenv`.
* **db**: `bun`.
* **redis**: `go-redis/v9`.
* **swagger**: `huma`.
* **httpclient**: `resty`.
* **apm/tracing***: `otel` with `inatrace` (backend `datadog`, TODO `elastic`)

## Project Structures

* cmd: for cli commands
* config: configuration
* db/migrations: db migration steps
* internal/handler: handlers for API
* internal/model: model or table
* internal/service: app use case / business logic / services / repositories
* lib/cache: cache eg. redis
* lib/db: database eg. postres
* lib/transport: fiber and endpoints or routes

To enable connectors rename by ommiting `.off` suffix, then do `go mod tidy`.

## Development

Install mockery command to generate interfece mock

as go 1.22
```
$ go install github.com/vektra/mockery/v2@v2.46.0
```

To re-generate mock from interface:

```
$ mockery --all
```
