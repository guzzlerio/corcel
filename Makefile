all: install clean build lint test

clean:
	go clean

build: clean
	go build

test: build lint
	go test

lint:
	bash scripts/lint.sh

install:
	go get -t
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update

.PHONY: clean build lint test install
