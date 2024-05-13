# Gophermart

## Requirements for run 

Installed [goose](https://github.com/pressly/goose) for run migrations

## Requirements for development 

Installed: 
* [goose](https://github.com/pressly/goose) for run migrations
* [swag](https://github.com/swaggo/swag) for generate OpenAPI
* [jet](https://github.com/go-jet/jet) for generate type safe queries
* [buf](https://github.com/bufbuild/buf) for generate proto

## Run

1. Setup pg db
2. Up migrations
```bash
make migrations_up
```
3. Start accrual server

```bash
chmod +x ./cmd/accrual/accrual_linux_amd64
./cmd/accrual/accrual_linux_amd64
```

4. Start server
```bash
go run ./cmd/gophermart
```