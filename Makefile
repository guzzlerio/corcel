version=$(shell git -C $(CURDIR) describe --abbrev=0)
all: clean build lint test

clean:
	go clean

gen:
	if [ -d "corcel-reports-workbench" ]; then \
		(cd corcel-reports-workbench && git pull) \
	else \
		git clone git@github.com:guzzlerio/corcel-reports-workbench.git; \
	fi
	cd corcel-reports-workbench && npm install -d && npm install gulp && npm install -g gulp-cli && gulp && cp out/index.html ../report/data/corcel.layout.mustache.html


build: clean lint generate
	#version=`grep -Po "(?<=version=)[0-9.]+" version`
	#go build -ldflags="-X main.Version=${version}"
	go build -ldflags "-X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.CommitHash=`git rev-parse HEAD` -X main.Version=`git describe --tags --always`"

test: build 
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	ginkgo -cover -r --race -noisyPendings=false -slowSpecThreshold=10

generate:
	go get -u github.com/jteeuwen/go-bindata/...
	(cd report && go-bindata -pkg report data)

lint: generate
	go get -u github.com/alecthomas/gometalinter
	go get -t ./...
	gometalinter --install
	# The Before and AfterTest ones here are because deadcode in the linter collection falsely detects them as dead code
	gometalinter -e "BeforeTest" -e "AfterTest" -e "should have comment" -e "duplicate" -e "undeclared name: Asset" -e "_test" --deadline=600s

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
	bash ./scripts/demo.sh


.PHONY: clean build lint test install ui dist dist_linux demo generate
