.PHONY: all install-tools lint lint-fix clean help

LINT_VERSION := v2.13.2
BIN_DIR := $(CURDIR)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

all: lint ## Run the linter (default)

$(GOLANGCI_LINT):
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINT_VERSION)

install-tools: $(GOLANGCI_LINT) ## Install all development tools

lint: $(GOLANGCI_LINT) ## Run the linter
	$(GOLANGCI_LINT) run

lint-fix: $(GOLANGCI_LINT) ## Run the linter and auto-fix issues
	$(GOLANGCI_LINT) run --fix

clean: ## Remove downloaded tools
	rm -rf $(BIN_DIR)

help: ## Show available commands
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk -F ':.*?## ' '{printf "  %-16s %s\n", $$1, $$2}' | sort
