# Automatically detected modules with go.mod (including nested modules)
MODULES := $(shell find . -name go.mod -not -path "*/\.*" -exec dirname {} \; | grep -v '^\.$$' | sed 's|^\./||' | sort)

# Module categories
ESSENTIAL_MODULES := $(shell find essentials -name go.mod -exec dirname {} \; | sed 's|^\./||' | sort)
LANGUAGE_MODULES := $(shell if [ -d languages ]; then find languages -name go.mod -exec dirname {} \; | sed 's|^\./||' | sort; fi)
PIPELINE_MODULES := $(shell find pipelines -name go.mod -exec dirname {} \; | sed 's|^\./||' | sort)

# Colors for pretty output
GREEN  := \033[32m
YELLOW := \033[33m
WHITE  := \033[37m
RESET  := \033[0m
RED    := \033[31m

# Separator for visual clarity
define print_separator
	@echo "${WHITE}"; \
	printf '=%.0s' {1..80}; \
	echo "${RESET}"; \
	echo "${GREEN}$$1${RESET}"; \
	echo "${WHITE}"; \
	printf '=%.0s' {1..80}; \
	echo "${RESET}"
endef

# Variables to store development results
DEVELOP_SUCCESS := 
DEVELOP_FAILED := 

# Function to update development results
define update_develop_results
	@status=$$?; \
	if [ $$status -eq 0 ]; then \
		DEVELOP_SUCCESS="$$DEVELOP_SUCCESS$$1 "; \
	else \
		DEVELOP_FAILED="$$DEVELOP_FAILED$$1 "; \
	fi
endef

# Function to print development summary
define print_develop_summary
	@echo "\n${WHITE}=== Development Summary ===${RESET}"
	@if [ -n "$$DEVELOP_SUCCESS" ]; then \
		echo "${GREEN}Successful modules:${RESET}"; \
		for module in $$DEVELOP_SUCCESS; do \
			echo "  ✓ $$module"; \
		done; \
	fi
	@if [ -n "$$DEVELOP_FAILED" ]; then \
		echo "\n${RED}Failed modules:${RESET}"; \
		for module in $$DEVELOP_FAILED; do \
			echo "  ✗ $$module"; \
		done; \
		exit 1; \
	fi
	@if [ -z "$$DEVELOP_SUCCESS" ] && [ -z "$$DEVELOP_FAILED" ]; then \
		echo "${YELLOW}No modules were processed${RESET}"; \
	fi
endef

# Default target
.DEFAULT_GOAL := help

# Ensure dependencies are installed
check-deps:
	@command -v go >/dev/null 2>&1 || { echo "${YELLOW}go is not installed. Please install it.${RESET}"; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "${YELLOW}golangci-lint is not installed. Please install it.${RESET}"; exit 1; }
	@command -v dagger >/dev/null 2>&1 || { echo "${YELLOW}dagger is not installed. Please install it.${RESET}"; exit 1; }

# Check if MODULE is valid
check-module:
	@if [ -z "$(MODULE)" ]; then \
		echo "${YELLOW}Error:${RESET} MODULE is required. Use MODULE=<module>"; \
		exit 1; \
	fi
	@if [ -n "$(MODULE)" ] && ! echo "$(MODULES)" | tr ' ' '\n' | grep -q "^$(MODULE)$$"; then \
		echo "${YELLOW}Error:${RESET} Invalid module: $(MODULE)"; \
		echo "Available modules by category:"; \
		echo "${GREEN}Essential modules:${RESET}"; \
		echo "$(ESSENTIAL_MODULES)" | tr ' ' '\n' | sed 's/^/  - /'; \
		echo "${GREEN}Language modules:${RESET}"; \
		echo "$(LANGUAGE_MODULES)" | tr ' ' '\n' | sed 's/^/  - /'; \
		echo "${GREEN}Pipeline modules:${RESET}"; \
		echo "$(PIPELINE_MODULES)" | tr ' ' '\n' | sed 's/^/  - /'; \
		exit 1; \
	fi

.PHONY: help
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  make <target> [MODULE=<module>]'
	@echo ''
	@echo 'Available modules by category:'
	@echo "${GREEN}Essential modules:${RESET}"
	@echo "$(ESSENTIAL_MODULES)" | tr ' ' '\n' | sed 's/^/  - /'
	@echo "${GREEN}Language modules:${RESET}"
	@echo "$(LANGUAGE_MODULES)" | tr ' ' '\n' | sed 's/^/  - /'
	@echo "${GREEN}Pipeline modules:${RESET}"
	@echo "$(PIPELINE_MODULES)" | tr ' ' '\n' | sed 's/^/  - /'
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
			$(call print_separator,"Cleaning $$module"); \
			(cd $$module && rm -rf internal dagger.gen.go go.sum 2>/dev/null || true); \
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
			$(call print_separator,"Tidying $$module"); \
			(cd $$module && go mod tidy) || exit 1; \
		done; \
	fi

.PHONY: develop
develop: check-deps ## Run dagger develop (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Developing $(MODULE)...${RESET}"; \
		(cd $(MODULE) && dagger develop); \
		$(call update_develop_results,$(MODULE)); \
	else \
		echo "${GREEN}Running dagger develop on all modules...${RESET}"; \
		for module in $(MODULES); do \
			$(call print_separator,"Developing $$module"); \
			(cd $$module && dagger develop); \
			$(call update_develop_results,$$module); \
		done; \
	fi
	$(call print_develop_summary)

.PHONY: test
test: check-deps ## Run tests (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Testing $(MODULE)...${RESET}"; \
		(cd $(MODULE) && go test ./...); \
	else \
		echo "${GREEN}Running tests on all modules in parallel...${RESET}"; \
		for module in $(MODULES); do \
			$(call print_separator,"Testing $$module"); \
			(cd $$module && go test ./...) || exit 1; \
		done; \
	fi

.PHONY: fmt
fmt: check-deps ## Format Go files (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Formatting $(MODULE)...${RESET}"; \
		(cd $(MODULE) && go fmt ./...); \
	else \
		echo "${GREEN}Formatting Go files in all modules...${RESET}"; \
		for module in $(MODULES); do \
			$(call print_separator,"Formatting $$module"); \
			(cd $$module && go fmt ./...); \
		done; \
	fi

.PHONY: vet
vet: check-deps ## Run go vet (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Vetting $(MODULE)...${RESET}"; \
		(cd $(MODULE) && go vet ./...); \
	else \
		echo "${GREEN}Running go vet on all modules...${RESET}"; \
		for module in $(MODULES); do \
			$(call print_separator,"Vetting $$module"); \
			(cd $$module && go vet ./...); \
		done; \
	fi

.PHONY: lint
lint: check-deps ## Run linters (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Linting $(MODULE)...${RESET}"; \
		(cd $(MODULE) && golangci-lint run ./...); \
	else \
		echo "${GREEN}Running linters on all modules...${RESET}"; \
		for module in $(MODULES); do \
			$(call print_separator,"Linting $$module"); \
			(cd $$module && golangci-lint run ./...); \
		done; \
	fi

.PHONY: build
build: check-deps ## Build modules (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Building $(MODULE)...${RESET}"; \
		(cd $(MODULE) && go build ./...); \
	else \
		echo "${GREEN}Building all modules...${RESET}"; \
		for module in $(MODULES); do \
			$(call print_separator,"Building $$module"); \
			(cd $$module && go build ./...); \
		done; \
	fi

.PHONY: update-deps
update-deps: check-deps ## Update dependencies (MODULE=<module> for specific module)
	@if [ -n "$(MODULE)" ]; then \
		echo "${GREEN}Updating dependencies for $(MODULE)...${RESET}"; \
		(cd $(MODULE) && go get -u ./... && go mod tidy); \
	else \
		echo "${GREEN}Updating dependencies for all modules...${RESET}"; \
		for module in $(MODULES); do \
			$(call print_separator,"Updating dependencies for $$module"); \
			(cd $$module && go get -u ./... && go mod tidy); \
		done; \
	fi

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test) (MODULE=<module> for specific module)

.PHONY: clean-temp
clean-temp: ## Clean temporary files
	@echo "${GREEN}Cleaning temporary files...${RESET}"
	@find . -name '*.tmp' -delete

.PHONY: essential
essential: check-deps ## Run operation on essential modules only (make essential <operation>)
	@if [ -z "$(OPERATION)" ]; then \
		echo "${YELLOW}Error:${RESET} OPERATION is required. Example: make essential OPERATION=test"; \
		exit 1; \
	fi
	@echo "${GREEN}Running $(OPERATION) on essential modules...${RESET}"; \
	for module in $(ESSENTIAL_MODULES); do \
		$(call print_separator,"$$module - $(OPERATION)"); \
		(cd $$module && make $(OPERATION)) || exit 1; \
	done

.PHONY: languages
languages: check-deps ## Run operation on language modules only (make languages <operation>)
	@if [ -z "$(OPERATION)" ]; then \
		echo "${YELLOW}Error:${RESET} OPERATION is required. Example: make languages OPERATION=test"; \
		exit 1; \
	fi
	@echo "${GREEN}Running $(OPERATION) on language modules...${RESET}"; \
	for module in $(LANGUAGE_MODULES); do \
		$(call print_separator,"$$module - $(OPERATION)"); \
		(cd $$module && make $(OPERATION)) || exit 1; \
	done

.PHONY: pipelines
pipelines: check-deps ## Run operation on pipeline modules only (make pipelines <operation>)
	@if [ -z "$(OPERATION)" ]; then \
		echo "${YELLOW}Error:${RESET} OPERATION is required. Example: make pipelines OPERATION=test"; \
		exit 1; \
	fi
	@echo "${GREEN}Running $(OPERATION) on pipeline modules...${RESET}"; \
	for module in $(PIPELINE_MODULES); do \
		$(call print_separator,"$$module - $(OPERATION)"); \
		(cd $$module && make $(OPERATION)) || exit 1; \
	done
