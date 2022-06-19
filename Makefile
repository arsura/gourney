dev: 
	nodemon --exec APP_ENV=development go run cmd/main.go --signal SIGTERM

dev-docker:
	docker-compose -f ./development/docker-compose.yml up -d

dev-docker-infra:
	docker-compose -f ./development/docker-compose.infra.yml up -d

clean-docker:
	docker rm gourney && docker image rm gourney

test:
	go clean -testcache && go test -v ./...



.PHONY: migrate-up-dev migrate-down-dev