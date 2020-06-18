#!/bin/bash

ProjectName=""
ProjectPath=""

while getopts ":n:p" opt
do
    case $opt in
    n)
        ProjectName=$OPTARG;;
    p)
        ProjectPath=$OPTARG;;
    ?)
        echo "-n project name"
        echo "-p project path"
        exit 1;;
    esac
done

if [ -z "$ProjectName" ]; then
    echo "must input project name"
    exit 1;
fi

echo "project name: $ProjectName"

if [ -z "$ProjectPath" ]; then
    echo "use default project path: ./$ProjectName/"
    ProjectPath="./$ProjectName/"
fi

echo "project path: $ProjectPath"

mkdir "$ProjectPath"

echo "start copy src to $ProjectPath"
rsync -av --exclude $ProjectPath --exclude .git --exclude .gitignore --exclude .idea --exclude .DS_Store --exclude go-svr-template ./ $ProjectPath
cd $ProjectPath
rm -rf start_project.sh $ProjectName go.mod go.sum ./log/* dev.sh
go mod init $ProjectName
sed -i "" "s/go-svr-template/$ProjectName/g" `grep go-svr-template -rl .`
Key=`openssl rand -hex 16`
sed -i "" "s/GinAuthKeyContent123456789012345/$Key/g" main.go