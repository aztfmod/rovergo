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
	go test ./pkg/... -count 1 -v

clean: ## 🧹 Cleanup project
	rm -rf bin
	rm -rf dist
	rm -rf landingzones
	go mod tidy

image-dev: ## 📦 Build the devcontainer image for Rover dev work
	docker build .devcontainer --file .devcontainer/Dockerfile \
	--tag $(IMAGE_REG)/$(IMAGE_REPO)-dev:$(IMAGE_TAG) \
	--build-arg INSTALL_GO=true \
	--build-arg INSTALL_DOCKER=false \
	--build-arg ROVER_VERSION=

image: ## 📦 Build the devcontainer image for Rover end users
	docker build .devcontainer --file .devcontainer/Dockerfile \
	--tag $(IMAGE_REG)/$(IMAGE_REPO):$(IMAGE_TAG) \
	--build-arg INSTALL_GO=false \
	--build-arg INSTALL_DOCKER=false \
	--build-arg ROVER_VERSION=latest

push: ## 🔼 Push the devcontainer images
	docker push $(IMAGE_REG)/$(IMAGE_REPO):$(IMAGE_TAG)
	docker push $(IMAGE_REG)/$(IMAGE_REPO)-dev:$(IMAGE_TAG)