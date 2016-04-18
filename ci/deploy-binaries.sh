for i in `ls dist`
do
    inArr=(${i//_/ })
    os=${inArr[0]}
    arch=${inArr[1]}
    version="$(git describe --abbrev=0)"
    token=$CORCEL_DEPLOY_TOKEN
    echo "curl -X POST -H 'Content-Type: application/octet-stream' --data-binary '@data.output' 'https://corcel.io/releases?version=$version&arch=$arch&os=$os&token=$token'"
done
