all: clean build lint test

clean:
	go clean

build: clean
	go build

test: build lint
	go test

lint:
	bash scripts/lint.sh

.PHONY: clean build lint test
