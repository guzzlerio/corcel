if [ "$TRAVIS_BRANCH" == "master" ]
then
    gox -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"
	ghr --username reaandrew -c master --token $GITHUB_TOKEN --replace `git describe --abbrev=0 --tags` dist/  
fi

if [ "$TRAVIS_BRANCH" == "release" ]
then
    gox -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"
	ghr --username reaandrew -c release --token $GITHUB_TOKEN --delete --replace --prerelease pre-release dist/  
fi

if [ "$TRAVIS_BRANCH" == "develop" ]
then
    gox -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"
	ghr --username reaandrew -c develop --token $GITHUB_TOKEN --delete --replace --prerelease latest dist/  
fi
