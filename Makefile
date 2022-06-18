run-dev-docker:
	docker-compose -f ./development/docker-compose.yml up

clean-docker:
	docker rm gourney && docker image rm gourney

test:
	go clean -testcache && go test -v ./...

dev: 
	nodemon --exec APP_ENV=development go run cmd/main.go --signal SIGTERM

.PHONY: migrate-up-dev migrate-down-dev