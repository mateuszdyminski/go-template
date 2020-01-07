GOLANG_VERSION := 1.13.5
ALPINE_VERSION := 3.11

NAME ?= $(shell echo $${PWD\#\#*/})
VERSION ?= $(shell git describe --always)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
LAST_COMMIT_USER ?= $(shell git log -1 --format='%cn <%ce>')
LAST_COMMIT_HASH ?= $(shell git log -1 --format=%H)
LAST_COMMIT_TIME ?= $(shell git log -1 --format=%cd --date=format:'%Y-%m-%d %H:%M:%S')

GIT_REPO := github.com/mateuszdyminski/go-template
DOCKER_REPO := mateuszdyminski

GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)

.DEFAULT_GOAL := all
.PHONY: all lint test build docker-build docker-push release misspell vet fmt run race help swag deps

all: deps fmt vet lint test build ## Combines `fmt` `vet` `lint` `test` `build` commands

misspell: ## Runs misspell
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo 'Please install "misspell" tool: https://github.com/client9/misspell'; \
		exit 1; \
	fi
	misspell -w $(GOFILES)

lint: ## Runs golangci-lint
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo 'Please install "golangci-lint" tool: https://github.com/golangci/golangci-lint'; \
		exit 1; \
	fi
	golangci-lint run -v

test: ## Runs all unit tests and generates coverage report
	go test -coverprofile cover.out -v ./...
	go tool cover -html=cover.out -o coverage_report.html

race: ## Runs all unit tests with -race flag
	go test -race -v ./...

vet: ## Runs `go vet` for all packages
	go vet $(PACKAGES)

swag: ## Generates Swagger documentation in /swagger-docs directory
	@hash swag > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo 'Please install "swag" tool: https://github.com/swaggo/swag'; \
		exit 1; \
	fi
	swag init -o=swagger-docs

fmt: ## Formats all Golang files
	gofmt -l -s -w $(GOFILES)

build: ## Builds App locally
	CGO_ENABLED=0 \
	go build \
	-v \
	-ldflags "-s -w \
	-X '$(GIT_REPO)/api.AppName=$(NAME)' \
	-X '$(GIT_REPO)/api.AppVersion=$(VERSION)' \
	-X '$(GIT_REPO)/api.APIVersion=$(VERSION)' \
	-X '$(GIT_REPO)/api.BuildTime=$(BUILD_TIME)' \
	-X '$(GIT_REPO)/api.LastCommitUser=$(LAST_COMMIT_USER)' \
	-X '$(GIT_REPO)/api.LastCommitHash=$(LAST_COMMIT_HASH)' \
	-X '$(GIT_REPO)/api.LastCommitTime=$(LAST_COMMIT_TIME)'" \
	-o $(NAME)-$(VERSION) .

docker-build: ## Builds Docker image with App
	docker build --pull \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg NAME="$(NAME)" \
	--build-arg VERSION="$(VERSION)" \
	--build-arg BUILD_TIME="$(BUILD_TIME)" \
	--build-arg LAST_COMMIT_USER="$(LAST_COMMIT_USER)" \
	--build-arg LAST_COMMIT_HASH="$(LAST_COMMIT_HASH)" \
	--build-arg LAST_COMMIT_TIME="$(LAST_COMMIT_TIME)" \
	--label="build.version=$(VERSION)" \
	--label="build.time=$(BUILD_TIME)" \
	--label="commit.user=$(LAST_COMMIT_USER)" \
	--label="commit.hash=$(LAST_COMMIT_HASH)" \
	--label="commit.time=$(LAST_COMMIT_TIME)" \
	--tag="$(DOCKER_REPO)/$(NAME):latest" \
	--tag="$(DOCKER_REPO)/$(NAME):$(VERSION)" \
	.

docker-push: ## Pushes current build version of Docker image to the registry
	docker push "$(DOCKER_REPO)/$(NAME):latest"
	docker push "$(DOCKER_REPO)/$(NAME):$(VERSION)"

release: deps docker-build docker-push ## Combines `deps`, `docker-build` and `docker-push` commands

deps: ## Synchronises all dependencies
	go mod tidy

run: swag ## Runs App in development mode locally
	DEBUG="true" \
	APP_HTTP_PORT="8080" \
	APP_HTTP_PPROF_PORT="8090" \
	APP_HTTP_GRACEFUL_TIMEOUT="10" \
	APP_HTTP_GRACEFUL_SLEEP="0"  \
	APP_POSTGRES_HOST="localhost" \
	APP_POSTGRES_PORT=5432 \
	APP_POSTGRES_USER="postgres" \
	APP_POSTGRES_PASSWORD="password" \
	APP_POSTGRES_DBNAME="app_db" \
	go run .

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
