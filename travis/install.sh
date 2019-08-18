#!/bin/bash
set -euC
set -o xtrace

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure && dep status

# Only run GopherJS tests and Linter tests on Linux/Go 1.12
if [[ "${TRAVIS_GO_VERSION}" == "1.12.x" ]]; then
    # Install nodejs and dependencies, but only for Linux
    curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
    sudo apt-get update -qq
    sudo apt-get install -y nodejs
    npm install
    # Then install GopherJS and related dependencies
    go get -u github.com/gopherjs/gopherjs

    # Source maps (mainly to make GopherJS quieter; I don't really care
    # about source maps in Travis)
    npm install source-map-support

    # Set up GopherJS for syscalls
    (
        cd $GOPATH/src/github.com/gopherjs/gopherjs/node-syscall/
        npm install --global node-gyp
        node-gyp rebuild
        mkdir -p ~/.node_libraries/
        cp build/Release/syscall.node ~/.node_libraries/syscall.node
    )

    go get -u -d -tags=js github.com/gopherjs/jsbuiltin

    # Linter
    curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.17.1
fi
