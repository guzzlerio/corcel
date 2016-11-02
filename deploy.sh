branch=`git rev-parse --abbrev-ref HEAD`

token=$GO_DEPLOY_GITHUB_TOKEN
user=$GO_DEPLOY_GITHUB_USER
repo=$GO_DEPLOY_GITHUB_REPO

tmpDirectory="/tmp/githubpub/$(uuidgen)"
mkdir -p $tmpDirectory

trap "rm -rf $tmpDirectory" EXIT

if ["$branch" == "develop"]
then
    echo "Deploying the latest development edition"

    # Determine if there is a release called "latest"
    curl -u $user:$token "https://api.github.com/repos/$user/$repo/releases" > "$tmpDirectory/releases.json"

    echo "$tmpDirectory/releases.json"
fi

if ["$branch" == "release"]
then
    echo "Deploying the next as a pre-release"
fi

if ["$branch" == "master"]
then
    echo "Deploying the latest stable"
fi
