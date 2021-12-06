migrate-up-dev:
	migrate -path pkg/models/pgsql/migrations -database "postgres://admin:admin@localhost:8081/gourney?sslmode=disable" -verbose up

migrate-down-dev:
	migrate -path pkg/models/pgsql/migrations -database "postgres://admin:admin@localhost:8081/gourney?sslmode=disable" -verbose down

dev:
	docker-compose -f ./development/docker-compose.yml up

test:
	go test -v ./...

.PHONY: migrate-up-dev migrate-down-dev