GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME ?= $(shell basename $(PWD))## The name of the binary.  Defaults to the name of the current directory
GOOS ?= $(shell $(GOCMD) env GOOS)## The operating system to build for.  Defaults to the current operating system
GOARCH ?= $(shell $(GOCMD) env GOARCH)## The architecture to build for.  Defaults to the current architecture
BUILD_ROOT ?= build/## The root directory for build artifacts.  Defaults to build/
BIN_PATH ?= $(BUILD_ROOT)/$(BINARY_NAME)## The path to the binary.  Defaults to $BUILD_ROOT/$BINARY_NAME
VERSION ?= 0.0.1## The version of the service.  Defaults to 0.0.1
APP_ENV ?= local## The environment to deploy to.  Defaults to local
DIR := $(shell pwd -L)
SERVICE_NAME ?= $(shell basename $(PWD))## The name of the service.  Defaults to the name of the current directory

.PHONY: test build clean tidy run help

help: ## Display this help message
	@echo "Usage:"
	@echo "  make <target> <OPTIONAL_ENV_VAR=value>"
	@echo ""
	@echo "Targets:"
	@echo "------------------------------------------------------------------"
	@max=$$(awk 'BEGIN { max = 0 } /^[a-zA-Z_-]+:.*?## .*$$/ { if (length($$1) > max) max = length($$1) } END { print max }' $(MAKEFILE_LIST)); \
    grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
    sort | \
    awk -v len=$$max 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-*s\033[0m %s\n", len+4, $$1, $$2}'
	@echo "------------------------------------------------------------------"
	@echo ""
	@echo "Environment Variables:"
	@echo "------------------------------------------------------------------"
	@max=$$(awk 'BEGIN { max = 0 } /^[a-zA-Z_]+ ?\?=.*?## .*$$/ { if (length($$1) > max) max = length($$1) } END { print max }' $(MAKEFILE_LIST)); \
	grep -E '^[a-zA-Z_]+ ?\?=.*?## .*$$' $(MAKEFILE_LIST) | sort | awk -v len=$$max 'BEGIN {FS = " ?\?=.*?## "}; {printf "\033[33m%-*s\033[0m %s\n", len+4, $$1, $$2}'
	@echo "------------------------------------------------------------------"

# TODO: Fix all this junk below.

test: ## Run any tests in the src directory on your local machine
	@echo "Running tests..."
	@cd src && $(GOCMD) clean -testcache && $(GOTEST) -v ./
	@cd ..

build: tidy test## Build the project
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build -C src -a -o $(BIN_PATH) -ldflags='-s -w -extldflags "-static"' -trimpath
# @chmod +x ./scripts/bash/*
# @./scripts/bash/build.sh --version $(VERSION)
clean: ## Clean up build objects
	@echo "Cleaning up..."
	@cd src && $(GOCMD) clean -testcache && rm -rf build/* && rm -rf tmp/* && rm -rf service/* && rm -rf service/.env
	@cd ..
	@rm -rf build/* && rm -rf tmp/* && rm -rf service/* && rm -rf service/.env

tidy: ## Run go tidy
	@echo "Tidying up..."
	@$(GOCMD) mod tidy

run: env hotreload ## Run the service locally using hotreload

ignore: ## Add build/ to .gitignore
# check if there is a gitignore file, if not create one
	@if [ ! -f .gitignore ]; then \
		touch .gitignore; \
	fi
#check if .gitignore has "src/build/" in it, if not add it
	@if ! grep -q "src/build" .gitignore; then \
		echo "src/build" >> .gitignore; \
	fi
# check if .gitignore has "service" in it, if not add it
	@if ! grep -q "service" .gitignore; then \
		echo "service" >> .gitignore; \
	fi
# check if .gitignore has "src/tmp" in it, if not add it
	@if ! grep -q "tmp" .gitignore; then \
		echo "tmp" >> .gitignore; \
	fi

env: ## Generate the src/service/.env file for local development
# check if the service directory exists, if not create it
	@[[ -d service ]] || { mkdir -p service; }
# check if the .env file exists, if not create it, if it does, remove it and create a new one
	@[[ -f service/.env ]] && { rm service/.env; touch service/.env; } || { touch service/.env; }
ifeq ($(shell [ -f .env_secrets ] && echo true),true)
	@for i in `cat .env_secrets` ; do echo "$$i" >> service/.env ; done
endif

localtest: ## Run any tests in the src/tests directory on your local machine
	@cd src && $(GOCMD) clean -testcache && DOCKER_GATEWAY_HOST=localhost $(GOCMD) test -v ./tests/...
	@cd ..

testall: test localtest ## Run all tests

local: env hotreload ## Run the service locally using hotreload

hotreload: ## Run docker-compose-hotreload.yaml to run the service locally using Air
	@docker compose -f docker-compose-hotreload.yaml up

build-local-docker: ## Build local dockerfile for hotreload
	@echo "Building $(SERVICE_NAME) local container"
# implement local build stuff.

status: ## Get the status of the service from 
	@curl -s localhost:8080/status | jq

healthcheck: ## Make a request to the /healthcheck endpoint
	@curl -s localhost:8080/healthcheck | jq