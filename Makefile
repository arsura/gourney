migrate-up-dev:
	migrate -path pkg/models/pgsql/migrations -database "postgres://admin:admin@localhost:8091/gourney?sslmode=disable" -verbose up

migrate-down-dev:
	migrate -path pkg/models/pgsql/migrations -database "postgres://admin:admin@localhost:8091/gourney?sslmode=disable" -verbose down

run-dev-docker:
	docker-compose -f ./development/docker-compose.yml up

clean-docker:
	docker rm gourney && docker image rm gourney

test:
	go clean -testcache && go test -v ./...

.PHONY: migrate-up-dev migrate-down-dev