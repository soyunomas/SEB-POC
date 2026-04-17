# SEB-Linux Strict Makefile
.PHONY: all build test lint clean run

export CGO_ENABLED=0

all: lint test build

build:
	@echo "=> Building SEB-Linux (CGO_ENABLED=0)..."
	go build -ldflags="-s -w" -o bin/seb-linux ./cmd/seb

test:
	@echo "=> Running tests with race detector..."
	CGO_ENABLED=1 go test -race -v ./...

lint:
	@echo "=> Running linters & CGO checks..."
	go vet ./...
	@go list -f '{{.ImportPath}}: {{.CgoFiles}}' ./... | awk '/\[\]$$/ {next} {print "[ERROR] CGO DETECTADO en " $$0; exit 1}'

clean:
	@echo "=> Cleaning up..."
	rm -rf bin/

run: build
	@echo "=> Ejecutando SEB-Linux..."
	./bin/seb-linux
