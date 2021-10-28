# Only used locally, real releases are versioned with goreleaser
VERSION ?= 2.0.0-dev

.PHONY: build help lint lint-fix run test clean image image-dev push
.DEFAULT_GOAL := help

# Dev container images location
IMAGE_REG ?= symphonydev.azurecr.io
IMAGE_REPO ?= rover2
IMAGE_TAG ?= latest

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint # golangci-lint is a nightmare to run in a pipeline

default: ci

help: ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## 👀 Lint & format, will not fix but sets exit code on error
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	$(GOLINT_PATH) run --modules-download-mode=mod ./...

lint-fix: ## 🌟 Lint & format, will try to fix errors and modify code
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	$(GOLINT_PATH) run --modules-download-mode=mod ./... --fix

build: ## 🔨 Build the rover binary and place into ./bin/
	go build -ldflags "-X github.com/aztfmod/rover/pkg/version.Value='$(VERSION)'" -o bin/rover 

run: ## 🏃‍ Run locally, with hot reload, it's not very useful
	go run main.go $(ARGS)

test: ## 🤡 Run unit tests
	go test ./... -v  -tags unit

clean: ## 🧹 Cleanup project
	rm -rf bin
	rm -rf dist
	rm -rf landingzones
	go mod tidy

github:
	@bash "$(CURDIR)/scripts/build_image.sh" "github"

local:
	@bash "$(CURDIR)/scripts/build_image.sh" "local"

dev:
	@bash "$(CURDIR)/scripts/build_image.sh" "dev"

ci:
	@bash "$(CURDIR)/scripts/build_image.sh" "ci"

alpha:
	@bash "$(CURDIR)/scripts/build_image.sh" "alpha"
