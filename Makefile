GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=aoc

.PHONY: test clean tidy run new del

test: 
	@$(GOTEST) github.com/Zate/go-templates/...
clean: 
	@$(GOCMD) clean
	@rm -f $(BINARY_NAME)
tidy:
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
run:
	@$(GOBUILD) -o $(BINARY_NAME) -v .
	@./$(BINARY_NAME)
new:
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
del:
ifndef dir
	$(error existName is undefined)
endif
	@if [ -d $(dir) ]; then \
        rm -rf $(dir); \
        $(GOCMD) work edit -dropuse=$(dir); \
        echo "Deleted $(dir)"; \
    fi
