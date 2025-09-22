# SiGolang Boilerplate

* **mux**: `fiber` over `humafiber`.
* **logging**: `inalog`.
* **cli**: `spf13/cobra` over `humacli`.
* **config**: `cleanenv`.
* **db**: `bun`.
* **redis**: `go-redis/v9`.
* **swagger**: `huma`.
* **httpclient**: `resty`.
* **apm/tracing***: `otel` with `inatrace` (backend `uptrace`)

## Project Structures

* cmd: for cli commandsS
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

as go 1.25.1
```
$ go install github.com/vektra/mockery/v3@v3.5.4
```

To re-generate mock from interface:

```
$ mockery --all
```
