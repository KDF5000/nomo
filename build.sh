#!/bin/bash
set -e

ARCH=$1
ENV=$2
rm -rf output
mkdir -p output/bin output/conf output/static output/static/template

GIT_SHA=`git rev-parse --short HEAD || echo "NotGitVersion"`
GIT_BRANCH=`git branch 2>/dev/null | grep "^\*" | sed -e "s/^\*\ //"`
if [ "${GIT_BRANCH}" != "" ]; then
    GIT_SHA=$GIT_SHA"($GIT_BRANCH)"
fi
WHEN=`date '+%Y-%m-%d_%H:%M:%S'`

if [ "${ARCH}" == "linux" ];then
  export CGO_ENABLED=0
  export GOOS=linux
  export GOARCH=amd64
fi

env_path=./cmd/nomo/.env
if [  X"$ENV" != "X" ];
then
    env_path=./cmd/nomo/.env_${ENV}
fi

go build -installsuffix -a -v -o nomo -ldflags "-s -X main.GitSHA=${GIT_SHA} -X main.BuildTime=${WHEN}" ./cmd/nomo/

mv nomo output/bin/nomo
cp ${env_path} output/bin/.env
cp deploy/run.sh output/run.sh
cp deploy/run_wx.sh output/run_wx.sh
#cp deploy/conf/nomo.openhex.cn.crt output/conf/openhex.crt
#cp deploy/conf/nomo.openhex.cn.key output/conf/openhex.key
cp static/template/*.tpl output/static/template
