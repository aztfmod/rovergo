VERSION ?= 0.0.1

.PHONY: build
.DEFAULT_GOAL := help

help:  ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build:  ## 🔨 Build the rover binary
	go build -ldflags "-X main.version='$(VERSION)'" -o bin/rover 

run:  ## 🏃‍ Run locally
	go run main.go

clean:  ## 🧹 Cleanup project
	rm -rf bin