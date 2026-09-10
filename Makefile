.PHONY: all install-tools lint lint-fix coverage clean help

LINT_VERSION := v2.13.2
GOVERALLS_VERSION := v0.0.12
BIN_DIR := $(CURDIR)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
GOVERALLS := $(BIN_DIR)/goveralls

all: lint ## Run the linter (default)

$(GOLANGCI_LINT):
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINT_VERSION)

$(GOVERALLS):
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/mattn/goveralls@$(GOVERALLS_VERSION)

install-tools: $(GOLANGCI_LINT) $(GOVERALLS) ## Install all development tools

lint: $(GOLANGCI_LINT) ## Run the linter
	$(GOLANGCI_LINT) run

lint-fix: $(GOLANGCI_LINT) ## Run the linter and auto-fix issues
	$(GOLANGCI_LINT) run --fix

coverage: ## Generate the cross-package coverage profile exactly as CI does, then print it (goveralls sent to Coveralls by CI)
	go test -race -covermode atomic -coverprofile=covprofile -coverpkg ./... ./...
	@go tool cover -func=covprofile | tail -1

clean: ## Remove downloaded tools
	rm -rf $(BIN_DIR)

help: ## Show available commands
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk -F ':.*?## ' '{printf "  %-16s %s\n", $$1, $$2}' | sort
