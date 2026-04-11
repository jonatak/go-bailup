BINARY ?= bailup

.PHONY: all build test fmt vet clean

all: build

build:
	mkdir -p bin
	go build -o bin/$(BINARY) ./...

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -rf bin
