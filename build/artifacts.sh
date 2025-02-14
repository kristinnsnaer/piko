#!/bin/bash
set -euo pipefail

# Set the VERSION to $1, otherwise get it from `git describe`
GIT_VERSION=$(git describe || echo "NONE")
VERSION="${1:-$GIT_VERSION}"

declare -a arr=(
	"linux/amd64"
	"linux/arm64"
	"darwin/amd64"
	"darwin/arm64"
	"windows/amd64"
)
declare -a marr=(
	"server"
	"client"
)

mkdir -p bin/artifacts

for i in "${arr[@]}"
do
	GOOSARCH=$i
	GOOS=${GOOSARCH%/*}
	GOARCH=${GOOSARCH#*/}
	for j in "${marr[@]}"
	do
		BINARY_NAME=piko-$GOOS-$GOARCH-$j

		echo "Building $BINARY_NAME $VERSION..."
		GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-X github.com/andydunstall/piko/pkg/build.Version=$VERSION -X github.com/andydunstall/piko/pkg/build.Module=$j" -o bin/artifacts/$BINARY_NAME main.go
	done
done
