include .env

migrations_up:
	goose up

migrations_down:
	goose down

doc_gen:
	swag init -g "cmd/gophermart/main.go"

proto_gen:
	buf generate proto