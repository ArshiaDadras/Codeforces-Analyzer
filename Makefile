.PHONY: help test run build clean lint

help:
	@echo "Please use 'make <target>' where <target> is one of:"
	@echo "  test         to run tests"
	@echo "  run          to run the application"
	@echo "  build        to build the application"
	@echo "  clean        to remove the binary file"
	@echo "  lint         to run linter"
	@echo "  docker-build to build the docker image"
	@echo "  docker-run   to run the docker image"

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

docker-build:
	docker build -t codeforces-analyzer .
docker-run:
	docker run -p 8080:8080 codeforces-analyzer