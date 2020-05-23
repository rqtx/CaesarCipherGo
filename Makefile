# Get the latest tag
TAG=$(shell git describe --tags --abbrev=0)
GIT_COMMIT=$(shell git log -1 --format=%h)
GO_VERSION=1.14

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

go-run: ## Run program
	docker run --rm -v $$PWD/src:/app -v $$HOME/.ssh:/root/.ssh -w /app golang:$(GO_VERSION) go run main.go

go-build: ## Build binary
	docker run --rm -v $$PWD/src:/app -v $$HOME/.ssh:/root/.ssh -w /app golang:$(GO_VERSION) go build main.go

go-bash: ## Run bash to troubleshooting
	docker run --rm -i -v $$PWD/src:/app -v $$HOME/.ssh:/root/.ssh -w /app --entrypoint "bash" golang:$(GO_VERSION)
