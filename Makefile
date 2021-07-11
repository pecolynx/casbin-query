.PHONY: dependency unit-test docker-up docker-down

dependency:
	@go get -v ./...

unit-test: dependency
	@go test -v -short ./...

docker-up:
	@docker-compose -f docker/development/docker-compose.yml up -d
	sleep 20

docker-down:
	@docker-compose -f docker/development/docker-compose.yml down

test-docker-up:
	@docker-compose -f docker/test/docker-compose.yml up -d
	sleep 20

test-docker-down:
	@docker-compose -f docker/test/docker-compose.yml down
	