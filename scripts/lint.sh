MAX_COMPLEXITY=3

go get github.com/alecthomas/gometalinter
gometalinter --install --update

golint
gofmt -w -s ./*.go
goimports -w ./*.go
structcheck -t .
aligncheck .
go tool vet ./*.go
gocyclo -over ${MAX_COMPLEXITY} .
