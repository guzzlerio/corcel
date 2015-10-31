export GOPATH=/go
cd /go/src/bitbucket.org/guzzler/corcel
go get -t ./...
go build && go test
go get github.com/mitchellh/gox
gox -output "/var/ci/output/{{.OS}}_{{.Arch}}_{{.Dir}}"

