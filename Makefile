migrate-up:
	golang-migrate -path db/migration -database "postgres://admin:admin@localhost:5432/moonbase?sslmode=disable" -verbose up

migrate-down:
	golang-migrate -path db/migration -database "postgres://admin:admin@localhost:5432/moonbase?sslmode=disable" -verbose down


.PHONY: migrate-up migrate-down