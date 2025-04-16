# Variables
APP_NAME := capital-gains
DOCKER_COMPOSE := docker-compose

.PHONY: run-local
run-local:
	@echo "Running the application locally..."
	@go run cmd/main.go < input.txt

.PHONY: run-local2
run-local:
	@echo "Running the application locally..."
	@go run cmd/main.go < input.txt


.PHONY: build
build:
	@echo "Building the application binary..."
	@mkdir -p bin
	@go build -o bin/$(APP_NAME) cmd/main.go

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: run
run:
	@echo "Building image and running application via Docker Compose with piped input..."
	@$(DOCKER_COMPOSE) build
	@cat input.txt | $(DOCKER_COMPOSE) run --rm -T $(APP_NAME)

.PHONY: run2
run2:
	@echo "Building image and running application via Docker Compose with piped input..."
	@$(DOCKER_COMPOSE) build
	@cat input2.txt | $(DOCKER_COMPOSE) run --rm -T $(APP_NAME)

.PHONY: down
down:
	@echo "Stopping and removing application containers (if any)..."
	@$(DOCKER_COMPOSE) down --remove-orphans


.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run           - Run the application locally with input.txt"
	@echo "  make build         - Build the application binary"
	@echo "  make test          - Run all tests"
	@echo "  make clean         - Clean up build artifacts"
	@echo "  make up     		- Build and run the application via Docker, piping input.txt"
	@echo "  make down          - Stop and clean up Docker Compose containers"
