MAX_COMPLEXITY=3

golint
gofmt -w -s ./*.go
goimports -w ./*.go
structcheck -t .
aligncheck .
go tool vet ./*.go
gocyclo -over ${MAX_COMPLEXITY} .
errcheck .
go tool vet --shadow ./*.go
varcheck .
defercheck .
