BINARY ?= bailup

.PHONY: all build test test-short fmt vet clean install

all: build

install:
	go install ./cmd/$(BINARY)

build:
	mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/$(BINARY)

test:
	go test -v ./...

test-short:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -rf bin
