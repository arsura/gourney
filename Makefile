migrate-up:
	golang-migrate -path pkg/models/pgsql/migrations -database "postgres://admin:admin@localhost:5432/moonbase?sslmode=disable" -verbose up

migrate-down:
	golang-migrate -path pkg/models/pgsql/migrations -database "postgres://admin:admin@localhost:5432/moonbase?sslmode=disable" -verbose down

test:
	go test ./...

.PHONY: migrate-up migrate-down