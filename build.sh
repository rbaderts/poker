#!/bin/sh -x

echo date = `date`

set -e


if [ "$1" == "prod" ]; then

    
    cd $GOPATH/src/github.com/rbaderts/poker
    mkdir -p release/linux

    env GOOS=linux GOARCH=amd64 go build -o  ./release/server cmd/server/main.go && mv ./release/server ./release/linux
    result=$?

    if [ "$result" != "0" ];then
        echo "go build failed"
        exit 1
    fi
    mkdir -p ./release/linux/store

    cp -r ./web ./release/linux
    cp -r ./migrations ./release/linux/
    cp -r ./ssl ./release/linux/

    cp poker.prod.env ./release/linux
    tar cvf poker.tar ./release

elif [ "$1" == "docker" ]; then

    cd $GOPATH/src/github.com/rbaderts/poker
    rm -rf build
    mkdir build
    GOOS=linux GOARCH=amd64 go build -o build/server.linux main.go 
    cp -r web build
    cp -r migrations build

    cp -r initdb build
#    cp wait-for-it.sh build/
#    cp docker-entrypoint.sh build/
    cp .dockerignore build/
    cp Dockerfile build

    cd build
    docker build --no-cache=true -t poker/poker .
    cd ..

else 
    . ./poker.env
    . ./secrets.env

    export GOOS=darwin
    export OGARCH=amd64

    cd $GOPATH/src/github.com/rbaderts/poker

    go build 
    GOOS=darwin GOARCH=amd64 go build -o server main.go 


    ./server -cpuprofile=perf.log

    result=$?

fi


