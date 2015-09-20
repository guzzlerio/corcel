if [ "$TRAVIS_BRANCH" == "master" ]
then
	ghr --username reaandrew --token $GITHUB_TOKEN --replace `git describe --abbrev=0 --tags` dist/  
fi

if [ "$TRAVIS_BRANCH" == "release" ]
then
	ghr --username reaandrew --token $GITHUB_TOKEN --replace --prerelease pre-release dist/  
fi

if [ "$TRAVIS_BRANCH" == "develop" ]
then
	ghr --username reaandrew --token $GITHUB_TOKEN --replace --prerelease latest dist/  
fi
