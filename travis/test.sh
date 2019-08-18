#!/bin/bash
set -euC

go test -race ./...

# Only run GopherJS tests, linter and coverage tests on Linux/Go 1.12
if [[ ${TRAVIS_OS_NAME} == "linux" && "${TRAVIS_GO_VERSION}" == "1.12.x" ]]; then
    gopherjs test ./...

    # Linter
    golangci-lint run ./...

    # Coverage
    go test -coverprofile=coverage.txt -covermode=set && bash <(curl -s https://codecov.io/bash)
fi
