version=$(shell git -C $(CURDIR) describe --abbrev=0)
all: clean build lint test

clean:
	go clean

build: clean
	#version=`grep -Po "(?<=version=)[0-9.]+" version`
	go build -ldflags="-X main.Version=${version}"

test: build lint
	ginkgo -cover -r --race -noisyPendings=false -slowSpecThreshold=10

lint:
	bash scripts/lint.sh

install:
	go get -t ./...
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	go get github.com/alecthomas/gometalinter
	go get -u github.com/jteeuwen/go-bindata/...
	gometalinter --install --update

dist: dist_linux

dist_linux:
	env GOOS=linux GOARCH=amd64 go build -o dist/corcel_linux_amd64
	tar -zcf dist/corcel_linux_amd64.tar.gz dist/corcel_linux_amd64

ui: install
	(cd ui && npm install -d && gulp)
	go-bindata -o ui.generated.go ui/public/...

.PHONY: clean build lint test install ui dist dist_linux
