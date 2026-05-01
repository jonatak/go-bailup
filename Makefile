BINARY ?= baillconnect-to-mqtt
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || echo development)
COMMIT_SHA ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_LDFLAGS ?= -s -w -X github.com/jonatak/baillconnect-to-mqtt/internal/config.Version=$(VERSION) -X github.com/jonatak/baillconnect-to-mqtt/internal/config.CommitSHA=$(COMMIT_SHA) -X github.com/jonatak/baillconnect-to-mqtt/internal/config.BuildTime=$(BUILD_TIME)

.PHONY: all build test test-short fmt vet clean install

all: build

install:
	go install -ldflags "$(GO_LDFLAGS)" ./cmd/$(BINARY)

build:
	mkdir -p bin
	go build -ldflags "$(GO_LDFLAGS)" -o bin/$(BINARY) ./cmd/$(BINARY)

test:
	go test -v ./...

test-short:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run

clean:
	rm -rf bin
