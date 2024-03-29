VERSION ?= $(shell git describe --tags --always --match=v* || echo v0)
BUILD := $(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X=main.version=$(VERSION) -X=main.build=$(BUILD)"
LINTERFLAGS=--enable-all --disable gochecknoinits --disable gochecknoglobals --disable goimports --disable wsl --out-format=tab --tests=false
BUILDFLAGS=$(LDFLAGS)
PROJECTNAME=calendar
GOEXE := $(shell go env GOEXE)
GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
BINSERVER=bin/server$(GOEXE)
BINCLIENT=bin/client$(GOEXE)
BINSENDER=bin/sender$(GOEXE)
BINNOTIFICATOR=bin/notificator$(GOEXE)
MODULESERVER=github.com/Brialius/calendar/cmd/server
MODULECLIENT=github.com/Brialius/calendar/cmd/client
MODULESENDER=github.com/Brialius/calendar/cmd/sender
MODULENOTIFICATOR=github.com/Brialius/calendar/cmd/notificator
IMPORT_PATH := /usr/local/include
LINT_PATH := ./bin/golangci-lint
LINT_PATH_WIN := golangci-lint
LINT_SETUP := curl -sfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s latest
IMPORT_PATH_WIN := c:\protobuf\include

ifneq ($(GOOS), windows)
	RACE = -race
	PWD := $(shell pwd)
endif

ifeq ($(GOOS), windows)
	IMPORT_PATH := $(IMPORT_PATH_WIN)
	LINT_PATH := $(LINT_PATH_WIN)
	PWD := $(shell echo %cd%)
	LINT_SETUP := go install github.com/golangci/golangci-lint/cmd/golangci-lint
endif

export

.PHONY: setup
setup: ## Install all the build and lint dependencies
	$(LINT_SETUP)
	go install github.com/golang/protobuf/protoc-gen-go
	go get ./...

.PHONY: test
test: ## Run all the tests
	go test -cover $(RACE) -v $(BUILDFLAGS) ./...

.PHONY: lint
lint: ## Run all the linters
	$(LINT_PATH) run $(LINTERFLAGS) cmd/client
	$(LINT_PATH) run $(LINTERFLAGS) cmd/notificator
	$(LINT_PATH) run $(LINTERFLAGS) cmd/server
	$(LINT_PATH) run $(LINTERFLAGS) cmd/sender

.PHONY: ci
ci: setup lint build ## Run all the tests and code checks

.PHONY: generate
generate:
	protoc --go_out=plugins=grpc:internal/grpc api/api.proto -I $(IMPORT_PATH) -I .

.PHONY: build
build: clean mod-refresh build-server build-sender build-notificator build-client

.PHONY: build-server
build-server: mod-refresh ## Build a version
	go build $(BUILDFLAGS) -o $(BINSERVER) $(MODULESERVER)

.PHONY: build-sender
build-sender: mod-refresh ## Build a version
	go build $(BUILDFLAGS) -o $(BINSENDER) $(MODULESENDER)

.PHONY: build-notificator
build-notificator: mod-refresh ## Build a version
	go build $(BUILDFLAGS) -o $(BINNOTIFICATOR) $(MODULENOTIFICATOR)

.PHONY: build-client
build-client: mod-refresh ## Build a version
	go build $(BUILDFLAGS) -o $(BINCLIENT) $(MODULECLIENT)

.PHONY: install
install: mod-refresh ## Install a binary
	go install $(BUILDFLAGS)

.PHONY: clean
clean: ## Remove temporary files
	go clean $(BUILDFLAGS) $(MODULESERVER)
	go clean $(BUILDFLAGS) $(MODULECLIENT)
	go clean $(BUILDFLAGS) $(MODULESENDER)
	go clean $(BUILDFLAGS) $(MODULENOTIFICATOR)

.PHONY: mod-refresh
mod-refresh: ## Refresh modules
	go mod tidy -v

.PHONY: version
version:
	@echo $(VERSION)-$(BUILD)

.PHONY: release
release:
	git tag $(ver)
	git push origin --tags

.PHONY: integration-tests
integration-tests:
	echo Sleeping to wait test environment...
	sleep 40s
	go test -v ./integration_tests

.PHONY: deploy
deploy:
	docker-compose up -d --build

.PHONY: undeploy
undeploy:
	docker-compose down

.PHONY: deploy-tests
deploy-tests:
	docker-compose -f ./docker-compose.test.yaml up  -d --build
	docker-compose -f ./docker-compose.test.yaml logs --follow integration_tests

.PHONY: undeploy-tests
undeploy-tests:
	docker-compose -f docker-compose.test.yaml down

.DEFAULT_GOAL := build
