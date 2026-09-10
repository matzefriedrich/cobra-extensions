.PHONY: all install-tools lint lint-fix coverage goveralls clean help

LINT_VERSION := v2.13.2
GOVERALLS_VERSION := v0.0.12
BIN_DIR := $(CURDIR)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
GOVERALLS := $(BIN_DIR)/goveralls
GOVERALLS_REPORT := goveralls-report.json

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

goveralls: $(GOVERALLS) ## Write the Coveralls JSON report (goveralls-report.json) from the CI-style coverage profile
	go test -race -covermode atomic -coverprofile=covprofile -coverpkg ./... ./...
	@$(GOVERALLS) -coverprofile=covprofile -service=test -debug -allowgitfetch=false 2>&1 | \
		awk '/Posting data: /{sub(/^.*Posting data: /, ""); p=1} p{print} /^}/{exit}' > $(GOVERALLS_REPORT)
	@echo "Coverage report written to $(GOVERALLS_REPORT)"

clean: ## Remove downloaded tools
	rm -rf $(BIN_DIR)

help: ## Show available commands
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk -F ':.*?## ' '{printf "  %-16s %s\n", $$1, $$2}' | sort
