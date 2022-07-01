dev-api: 
	nodemon --exec APP_ENV=development ENABLE_API=1 go run cmd/main.go --signal SIGTERM

dev-worker: 
	nodemon --exec APP_ENV=development ENABLE_WORKER=1 go run cmd/main.go --signal SIGTERM

dev-docker-infra:
	docker-compose -f ./developments/docker-compose.infra.yml up -d

test:
	go clean -testcache && go test -v ./...



.PHONY: migrate-up-dev migrate-down-dev