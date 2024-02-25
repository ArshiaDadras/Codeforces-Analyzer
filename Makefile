.PHONY: help test run build clean lint

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@echo "  test   to run the tests"
	@echo "  run    to run the application"
	@echo "  build  to build the application"
	@echo "  clean  to remove the build artifacts"
	@echo "  lint   to run the linter"

test:
	go test -v tests/*.go

run:
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

clean:
	rm -rf bin

lint:
	golangci-lint run