all: clean build lint test

clean:
	go clean

build: clean
	go build

test: build lint
	ginkgo -r -noisyPendings=false -slowSpecThreshold=10

lint:
	bash scripts/lint.sh

install:
	go get -t ./...
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	go get github.com/alecthomas/gometalinter
	go get -u github.com/jteeuwen/go-bindata/...
	gometalinter --install --update

ui: install
	(cd ui && npm install -d && gulp)
	go-bindata -o ui.generated.go ui/public/...

.PHONY: clean build lint test install ui
