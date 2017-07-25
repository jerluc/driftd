VERSION=`git describe --tags 2>/dev/null || echo "untagged"`
COMMITISH=`git describe --always 2>/dev/null`
OS := $(shell uname)

all: clean get-deps build

clean:
	rm -rf target

Linux-deps:
	echo

Darwin-deps:
	brew install iproute2mac
	brew install Caskroom/cask/tuntap

get-deps: $(OS)-deps
	go get -u github.com/op/go-logging
	go get -u github.com/jerluc/gobee
	go get -u github.com/jerluc/serial
	go get -u github.com/songgao/water
	go get -u golang.org/x/net/ipv6
	go get -u gopkg.in/alecthomas/kingpin.v2

build:
	go build -o ./target/driftd -ldflags="-X main.Version=${VERSION} -X main.Commitish=${COMMITISH}"

install: get-deps
	rm -f ${GOPATH}/bin/driftd
	go install -ldflags="-X main.Version=${VERSION} -X main.Commitish=${COMMITISH}"

.PHONY: clean get-deps build install all
