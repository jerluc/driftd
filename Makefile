VERSION=`git describe --tags 2>/dev/null || echo "untagged"`
COMMITISH=`git describe --always 2>/dev/null`

all: clean get-deps build

clean:
	rm -rf target

get-deps:
	go get github.com/op/go-logging
	go get github.com/jerluc/gobee
	go get github.com/jerluc/serial
	go get github.com/songgao/water
	go get golang.org/x/net/ipv6
	go get gopkg.in/alecthomas/kingpin.v2

build:
	go build -o ./target/riftd -ldflags="-X main.Version=${VERSION} -X main.Commitish=${COMMITISH}"

install:
	rm -f ${GOPATH}/bin/riftd
	go install -ldflags="-X main.Version=${VERSION} -X main.Commitish=${COMMITISH}"

.PHONY: clean get-deps build install all
