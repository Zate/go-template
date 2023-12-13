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

test: ## Run tests
	@$(GOTEST) ./...
build: ## Build the project
	@chmod +x ./scripts/bash/*
	@./scripts/bash/build.sh --version $(VERSION)
clean: ## Clean up build objects
	@$(GOCMD) clean
	@rm -rf build/
tidy: ## Run go tidy
	@$(GOCMD) mod tidy
run: tidy build ## Run tests, build and then run the project
	@./build/$(BINARY_NAME)
ignore: ## Add build/ to .gitignore
# check if there is a gitignore file, if not create one
	@if [ ! -f .gitignore ]; then \
		touch .gitignore; \
	fi
#check if .gitignore has "build/" in it, if not add it
	@if ! grep -q "build/" .gitignore; then \
		echo "build/" >> .gitignore; \
	fi

