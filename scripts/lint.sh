
# The following line should be uncommented once a baseline has been set for the code complexity and
# styling in the solution.  With this uncommented it means that any break to style or complexity
# will break the build (which is now in a convenient build file)

#set -e



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
