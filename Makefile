GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=$(shell basename $(CURDIR))

.PHONY: test clean tidy run new del help

help: ## Display this help message
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

test: ## Run tests
	@$(GOTEST) github.com/Zate/go-templates/...
clean: ## Clean up build objects
	@$(GOCMD) clean
	@rm -f $(BINARY_NAME)
tidy: ## Run go tidy in all templates
	@for d in $(shell ls -d */); do \
        d=$${d%?}; \
        $(GOCMD) work use $$d; \
        echo "Tidying $$d"; \
        cd $$d; \
        if [ -f go.mod ]; then \
            if grep -q "module github.com/Zate/go-templates/$$d" go.mod; then \
                $(GOCMD) mod tidy; \
            else \
                echo "go.mod exists but does not contain the correct module name"; \
                rm go.mod; \
                $(GOCMD) mod init github.com/Zate/go-templates/$$d; \
                $(GOCMD) mod tidy; \
            fi; \
        else \
            echo "go.mod does not exist"; \
            $(GOCMD) mod init github.com/Zate/go-templates/$$d; \
            $(GOCMD) mod tidy; \
        fi; \
        cd ..; \
        $(GOCMD) work sync; \
    done
run: ## Run the project
	@$(GOBUILD) -o $(BINARY_NAME) -v .
	@./$(BINARY_NAME)
new: ## Create a new template from base using dir=newName
ifndef dir
	$(error newName is undefined)
endif
	@cp -r base $(dir)
	@rm -rf $(dir)/go.mod $(dir)/go.sum
	@$(GOCMD) work edit -use=$(dir)
	@cd $(dir); \
	$(GOCMD) mod init github.com/Zate/go-templates/$(dir); \
	$(GOCMD) mod tidy; \
	cd ..;
del: ## Delete a template using dir=delName
ifndef dir
	$(error existName is undefined)
endif
	@if [ -d $(dir) ]; then \
        rm -rf $(dir); \
        $(GOCMD) work edit -dropuse=$(dir); \
        echo "Deleted $(dir)"; \
    fi
