for i in `ls dist`
do
    inArr=(${i//_/ })
    os=${inArr[0]}
    arch=${inArr[1]}
    version="$(git describe)"
    token=$CORCEL_DEPLOY_TOKEN
    curl --retry 5 -X POST -H "Content-Type: application/octet-stream" --data-binary "@$1/${i}" "https://corcel.io/releases?version=$version&arch=$arch&os=$os&token=$token&latest=$2"
done
