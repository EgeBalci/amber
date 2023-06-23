BUILD=go build
BUILD_FLAGS=-trimpath -buildvcs=false -ldflags="-extldflags=-static -s -w -X github.com/egebalci/amber/config.Version=$$(git log --pretty=format:'v1.0.%at-%h' -n 1)" 

default:
	${BUILD} ${BUILD_FLAGS} -o amber