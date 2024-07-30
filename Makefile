BIN := "./bin/anti-bruteforce"
BIN_CLI := "./bin/cli"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/anti-bruteforce
	go build -v -o $(BIN_CLI) -ldflags "$(LDFLAGS)" ./cmd/cli

run: build
	docker-compose -f deployments/docker-compose.yaml up --build -d
	ANTI_BRUTEFORCE_REDIS_URL=redis://localhost:6379 $(BIN)

test:
	go test -race -count 10 ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker-compose -f deployments/docker-compose.yaml up --build -d

down:
	docker-compose -f deployments/docker-compose.yaml down

restart: down up

integration-test:
	set -e ;\
	docker-compose -f ./deployments/test-docker-compose.yaml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f ./deployments/test-docker-compose.yaml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f ./deployments/test-docker-compose.yaml down ;\
	docker-compose -f ./deployments/test-docker-compose.yaml rm ;\
	exit $$test_status_code ;

.PHONY: build run test lint