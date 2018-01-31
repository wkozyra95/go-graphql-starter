#!/usr/bin/env bash

SCRIPT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

set -x

function ensureStartDB {
    docker inspect mongodb > /dev/null 2>&1
    EXISTS=$?
    
    if [ "$EXISTS" != "0" ] ; then
        echo "Create mongo docker"
        docker run -p 27017:27017 --name mongodb -d mongo
    fi
    echo "Start mongo docker"
    docker start mongodb
}

if [ "$1" = "check" ]; then
    cd $SCRIPT_PATH
    ./run.sh lint
    ./run.sh test
elif [ "$1" = "lint" ]; then
    cd $SCRIPT_PATH
    gometalinter --config=.gometalinter.json --deadline 1000s ./...

elif [ "$1" = "test" ]; then
    cd $SCRIPT_PATH
    go test ./...

elif [ "$1" = "run" ]; then
    cd $SCRIPT_PATH
    ensureStartDB
    export DB_URL=mongodb://localhost:27017/mongodb
    export BACKEND_PORT=3000
    go run main.go

elif [ "$1" = "run:dev" ]; then
    cd $SCRIPT_PATH
    ensureStartDB
    export DB_URL=mongodb://localhost:27017/mongodb
    export BACKEND_PORT=3001
    gin

elif [ "$1" = "setup" ]; then
    cd $SCRIPT_PATH
    go get -u github.com/golang/dep/cmd/dep
    go get -u github.com/alecthomas/gometalinter
    go get -u github.com/codegangsta/gin
    gometalinter --install --force
else
    echo "\
Usage: $0 COMMAND [ARGS]

Commands:
    check
    lint
    test

    run
    run:dev

    setup
    "
fi



