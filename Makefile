version=$(shell git -C $(CURDIR) describe --abbrev=0)
all: clean build lint test

clean:
	go clean

build: clean
	#version=`grep -Po "(?<=version=)[0-9.]+" version`
	(cd report && go-bindata -pkg report data/...)
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

run_5: build
	go build && ./corcel run --summary --progress bar --workers 1 --duration 5s --plan .resources/sample-plan.yml

run_15: build
	go build && ./corcel run --summary --progress bar --workers 1 --duration 15s --plan .resources/sample-plan.yml

run_30: build
	go build && ./corcel run --summary --progress bar --workers 1 --duration 30s --plan .resources/sample-plan.yml

demo: build
	corcel run --progress bar --summary --duration 10s --plan .resources/sample-plan.yml

.PHONY: clean build lint test install ui dist dist_linux demo
