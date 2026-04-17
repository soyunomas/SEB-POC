# SEB-Windows Strict Makefile
.PHONY: all build test lint clean run

export CGO_ENABLED=0

all: lint test build

build:
	@echo "=> Building SEB-Windows (CGO_ENABLED=0)..."
	go build -ldflags="-s -w" -o bin/seb-windows.exe ./cmd/seb

test:
	@echo "=> Running tests with race detector..."
	CGO_ENABLED=1 go test -race -v ./...

lint:
	@echo "=> Running linters..."
	go vet ./...

clean:
	@echo "=> Cleaning up..."
	rm -rf bin/

run: build
	@echo "=> Ejecutando SEB-Windows..."
	./bin/seb-windows.exe
