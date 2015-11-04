rm -rf output
mkdir output
docker run -v "$(pwd)/scripts:/var/ci/build/" -v $(pwd)/output:/var/ci/output/ -v $(pwd):/go/src/bitbucket.org/guzzler/corcel -t golang:1.5 bash -c "bash /var/ci/build/build-all.sh"
