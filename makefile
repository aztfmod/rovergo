VERSION ?= 0.0.1

.PHONY: build help lint lint-fix run test clean
.DEFAULT_GOAL := help

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint # Remove if not using Go

help: ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## 👀 Lint & format, will not fix but sets exit code on error
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	$(GOLINT_PATH) run --modules-download-mode=mod ./...

lint-fix: ## 🌟 Lint & format, will try to fix errors and modify code
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	$(GOLINT_PATH) run --modules-download-mode=mod ./... --fix

build: ## 🔨 Build the rover binary
	go build -ldflags "-X github.com/aztfmod/rover/pkg/version.Value='$(VERSION)'" -o bin/rover 

run: ## 🏃‍ Run locally
	go run main.go $(ARGS)

test: ## 🤡 Run tests
	go test ./pkg/... -count 1 -v

clean: ## 🧹 Cleanup project
	rm -rf bin
	rm -rf dist
	rm -rf landingzones
	go mod tidy