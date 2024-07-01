GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=$(shell basename $(PWD))
VERSION=0.0.1

.PHONY: test build clean tidy run help

help: ## Display this help message
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

# TODO: Fix all this junk below.

# test: ## Run tests
# 	@$(GOTEST) ./...
# build: ## Build the project
# 	@chmod +x ./scripts/bash/*
# 	@./scripts/bash/build.sh --version $(VERSION)
# clean: ## Clean up build objects
# 	@$(GOCMD) clean
# 	@rm -rf build/
# tidy: ## Run go tidy
# 	@$(GOCMD) mod tidy
# run: tidy build ## Run tests, build and then run the project
# 	@./build/$(BINARY_NAME)
# ignore: ## Add build/ to .gitignore
# # check if there is a gitignore file, if not create one
# 	@if [ ! -f .gitignore ]; then \
# 		touch .gitignore; \
# 	fi
# #check if .gitignore has "build/" in it, if not add it
# 	@if ! grep -q "build/" .gitignore; then \
# 		echo "build/" >> .gitignore; \
# 	fi



# .DEFAULT_GOAL := help

# .PHONY : build build-local-docker clean env help hotreload local localtest status

# # APP_ENV ?= local## The environment to deploy to.  Defaults to local
# # DIR := $(shell pwd -L)
# # TEAM ?= asecurityteam## The team to deploy to.  Defaults to asecurityteam
# SERVICE_NAME ?= $(shell cat project-descriptor.yml | grep name | sed -e 's/ //g'| cut -d":" -f2)## The name of the service.  Defaults to the name in project-descriptor.yml
# # APP_VERSION ?= $(shell cat $(SERVICE_NAME).sd.yml | grep buildNumber | cut -d" " -f2)## The version of the service.  Defaults to the buildNumber in the service descriptor yaml file
# # BITBUCKET_BUILD_NUMBER ?= $(APP_VERSION)## The build number of the service.  Defaults to the buildNumber in the service descriptor yaml file
# # NETRC ?= $(shell [ -f ~/.netrc ] && echo true || echo false)## Whether or not to use ~/.netrc for authentication.  Defaults to false
# # GOOS ?= linux## The operating system to build for.  Defaults to linux
# # GOARCH ?= amd64## The architecture to build for.  Defaults to amd64
# # COMMIT ?= $(shell git rev-parse --short=12 HEAD)## The commit hash of the current git commit.  Defaults to the short hash of the current commit
# # SOX ?= 0## Whether or not to deploy to the sox environment.  Defaults to false
# # REGISTRY ?= docker.atl-paas.net## The registry to push to.  Defaults to docker.atl-paas.net
# # REGISTRY_USER ?= $(shell cat ~/.netrc | grep login | cut -d" " -f2)## The username to use for the registry.  Defaults to the username in ~/.netrc
# # IMAGE_PATH ?= /$(TEAM)## The path to the image in the registry.  Defaults to /<TEAM> unless SOX is set to 1, then it defaults to /sox/<TEAM>
# # ifeq ($(NETRC),TRUE)
# # REGISTRY_PWD ?= $(shell cat ~/.netrc | grep password | cut -d" " -f2)## The password to use for the registry.  Defaults to the password in ~/.netrc
# # endif
# # ifeq ($(SOX),1)
# # IMAGE_PATH ?= /sox/$(TEAM)
# # endif
# # ARTIFACT ?= $(REGISTRY)$(IMAGE_PATH)/$(SERVICE_NAME)## The full path to the image in the registry.  Defaults to <REGISTRY><IMAGE_PATH>/<SERVICE_NAME>

# # SDCLI_VERSION ?= latest## The version of the sdcli to use.  Defaults to latest
# # SDCLI ?= docker run --rm -v "$(DIR):$(DIR)" -w "$(DIR)" asecurityteam/sdcli:$(SDCLI_VERSION)

# # VULN ?= VULN-1104043## The vulnerability to submit to the service.  Defaults to VULN-1104043 for testing

# # STG_URL ?= $(shell atlas micros service show --service $(SERVICE_NAME) -o json | jq -r '.environments.["stg-east"].url')




# help: ## Display this help message
# 	@echo "Usage:"
# 	@echo "  make <target> <OPTIONAL_ENV_VAR=value>"
# 	@echo ""
# 	@echo "Targets:"
# 	@echo "------------------------------------------------------------------"
# 	@max=$$(awk 'BEGIN { max = 0 } /^[a-zA-Z_-]+:.*?## .*$$/ { if (length($$1) > max) max = length($$1) } END { print max }' $(MAKEFILE_LIST)); \
#     grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
#     sort | \
#     awk -v len=$$max 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-*s\033[0m %s\n", len+4, $$1, $$2}'
# 	@echo "------------------------------------------------------------------"
# 	@echo ""
# 	@echo "Environment Variables:"
# 	@echo "------------------------------------------------------------------"
# 	@max=$$(awk 'BEGIN { max = 0 } /^[a-zA-Z_]+ ?\?=.*?## .*$$/ { if (length($$1) > max) max = length($$1) } END { print max }' $(MAKEFILE_LIST)); \
# 	grep -E '^[a-zA-Z_]+ ?\?=.*?## .*$$' $(MAKEFILE_LIST) | sort | awk -v len=$$max 'BEGIN {FS = " ?\?=.*?## "}; {printf "\033[33m%-*s\033[0m %s\n", len+4, $$1, $$2}'
# 	@echo "------------------------------------------------------------------"


# env: ## Generate the src/service.env file for local development
# 	@mkdir -p src/service
# 	@touch src/service/.env
	
# ifeq ($(shell [ -f .env_secrets ] && echo true),true)
# 	@for i in `cat .env_secrets` ; do echo "$$i" >> src/service/.env ; done
# endif

# localtest: ## Run any tests in the src/tests directory on your local machine
# 	@cd src && go clean -testcache && DOCKER_GATEWAY_HOST=localhost go test -v ./tests/...
# 	@cd ..

# test: ## Run any tests in the src directory on your local machine
# 	@cd src && go clean -testcache && go test -v ./
# 	@cd ..

# testall: test localtest ## Run all tests

# local: env hotreload ## Run the service locally using hotreload

# hotreload: ## Run docker-compose-hotreload.yaml to run the service locally using Air
# 	@docker compose -f docker-compose-hotreload.yaml up

# build-local-docker: ## Build local dockerfile for hotreload
# 	@echo "Building $(SERVICE_NAME) local container"
# # implement local build stuff.

# build: ## Build the service using the Dockerfile

# clean: ## Clean up any artifacts from building and deploying
# 	@rm -rf bin
# 	@rm -rf src/service/.env

# status: ## Get the status of the service from http://localhost:8080/status

# debug: ## Make a request to the /debug endpoint