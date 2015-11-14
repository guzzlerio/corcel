
# The following line should be uncommented once a baseline has been set for the code complexity and
# styling in the solution.  With this uncommented it means that any break to style or complexity
# will break the build (which is now in a convenient build file)

#set -e



MAX_COMPLEXITY=4

echo "golint"
golint
echo "gofmt -w -s ./*.go"
gofmt -w -s ./*.go
echo "goimports -w ./*.go"
goimports -w ./*.go
echo "structcheck -t ."
structcheck -t .
echo "aligncheck ."
aligncheck .
echo "go tool vet ./*.go"
go tool vet ./*.go
echo "gocyclo -over ${MAX_COMPLEXITY} ."
gocyclo -over ${MAX_COMPLEXITY} .
echo "errcheck ."
errcheck .
echo "go tool vet --shadow ./*.go"
go tool vet --shadow ./*.go
echo "varcheck ."
varcheck .
echo "defercheck ."
defercheck .
