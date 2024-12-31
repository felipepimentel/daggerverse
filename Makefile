# Automatically detected modules with go.mod
MODULES := $(shell find . -maxdepth 2 -name go.mod -exec dirname {} \; | grep -v '^\.$$' | sed 's|^\./||')

# Colors for pretty output
GREEN  := \033[32m
YELLOW := \033[33m
WHITE  := \033[37m
RESET  := \033[0m

# Default target
.DEFAULT_GOAL := help

# Ensure dependencies are installed
check-deps:
	@command -v go >/dev/null 2>&1 || { echo "${YELLOW}go is not installed. Please install it.${RESET}"; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "${YELLOW}golangci-lint is not installed. Please install it.${RESET}"; exit 1; }

# Check if MODULE is valid
check-module:
	@if [ -z "$(MODULE)" ]; then \
		echo "${YELLOW}Error:${RESET} MODULE is required. Use MODULE=<module>"; \
		exit 1; \
	fi
	@if [ -n "$(MODULE)" ] && ! echo "$(MODULES)" | grep -q "\b$(MODULE)\b"; then \
		echo "${YELLOW}Error:${RESET} Invalid module: $(MODULE). Available modules: $(MODULES)"; \
		exit 1; \
	fi

.PHONY: help
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  make <target> [MODULE=<module>]'
	@echo ''
	@echo 'Available modules: $(MODULES)'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: clean
clean: check-deps ## Clean generated files and dependencies (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Cleaning $(MODULE)...${RESET}"; \
		cd $(MODULE) && rm -rf internal dagger.gen.go go.sum 2>/dev/null || true; \
	else \
		echo "${GREEN}Cleaning all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Cleaning $$module...${RESET}"; \
			cd $$module && rm -rf internal dagger.gen.go go.sum 2>/dev/null || true; \
		done; \
	fi

.PHONY: tidy
tidy: check-deps ## Run go mod tidy (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Running go mod tidy on $(MODULE)...${RESET}"; \
		(cd $(MODULE) && go mod tidy); \
	else \
		echo "${GREEN}Running go mod tidy on all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Tidying $$module...${RESET}"; \
			(cd $$module && go mod tidy) || exit 1; \
		done; \
	fi

.PHONY: develop
develop: check-deps ## Run dagger develop (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Developing $(MODULE)...${RESET}"; \
		cd $(MODULE) && dagger develop; \
	else \
		echo "${GREEN}Running dagger develop on all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Developing $$module...${RESET}"; \
			cd $$module && dagger develop; \
		done; \
	fi

.PHONY: test
test: check-deps ## Run tests (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Testing $(MODULE)...${RESET}"; \
		cd $(MODULE) && go test ./...; \
	else \
		echo "${GREEN}Running tests on all modules in parallel...${RESET}"; \
		echo $(MODULES) | xargs -n 1 -P 4 -I {} bash -c 'echo "Testing {}"; cd {} && go test ./...'; \
	fi

.PHONY: fmt
fmt: check-deps ## Format Go files (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Formatting $(MODULE)...${RESET}"; \
		cd $(MODULE) && go fmt ./...; \
	else \
		echo "${GREEN}Formatting Go files in all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Formatting $$module...${RESET}"; \
			cd $$module && go fmt ./...; \
		done; \
	fi

.PHONY: vet
vet: check-deps ## Run go vet (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Vetting $(MODULE)...${RESET}"; \
		cd $(MODULE) && go vet ./...; \
	else \
		echo "${GREEN}Running go vet on all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Vetting $$module...${RESET}"; \
			cd $$module && go vet ./...; \
		done; \
	fi

.PHONY: lint
lint: check-deps ## Run linters (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Linting $(MODULE)...${RESET}"; \
		cd $(MODULE) && golangci-lint run ./...; \
	else \
		echo "${GREEN}Running linters on all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Linting $$module...${RESET}"; \
			cd $$module && golangci-lint run ./...; \
		done; \
	fi

.PHONY: build
build: check-deps ## Build modules (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Building $(MODULE)...${RESET}"; \
		cd $(MODULE) && go build ./...; \
	else \
		echo "${GREEN}Building all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Building $$module...${RESET}"; \
			cd $$module && go build ./...; \
		done; \
	fi

.PHONY: update-deps
update-deps: check-deps ## Update dependencies (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Updating dependencies for $(MODULE)...${RESET}"; \
		cd $(MODULE) && go get -u ./... && go mod tidy; \
	else \
		echo "${GREEN}Updating dependencies for all modules...${RESET}"; \
		for module in $(MODULES); do \
			echo "${YELLOW}Updating dependencies for $$module...${RESET}"; \
			cd $$module && go get -u ./... && go mod tidy; \
		done; \
	fi

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test) (MODULE=<module> for specific module)

.PHONY: clean-temp
clean-temp: ## Clean temporary files
	@echo "${GREEN}Cleaning temporary files...${RESET}"
	@find . -name '*.tmp' -delete
