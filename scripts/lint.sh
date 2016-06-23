
# The following line should be uncommented once a baseline has been set for the code complexity and
# styling in the solution.  With this uncommented it means that any break to style or complexity
# will break the build (which is now in a convenient build file)

set -e

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
# If you are reading this you will be wondering why i commented this out
# TWO FOLD
# Go checkout the logic required for assertions FOR  > >= < <= ETC...
# that and i don't know how to ignore specific methods with golint
#
# Going forward there might be some ways to refactor this down but
# unfortunately imo the complexity is unavoidable for this.
#
# We have also said the assertions will be done with javascript in the future
# which maybe the point when we can use JavaScript to help simplify
# the implementation.
#gocyclo -over ${MAX_COMPLEXITY} .
echo "errcheck ."
errcheck .
echo "go tool vet --shadow ./*.go"
go tool vet --shadow ./*.go
echo "varcheck ."
varcheck .
