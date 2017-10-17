CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-chatterbox/dispatcher
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-chatterbox/receiver
	cp dispatcher/*.go src/github.com/whosonfirst/go-whosonfirst-chatterbox/dispatcher/
	cp receiver/*.go src/github.com/whosonfirst/go-whosonfirst-chatterbox/receiver/
	cp *.go src/github.com/whosonfirst/go-whosonfirst-chatterbox
	if test ! -d src; then mkdir src; fi
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/eltorocorp/cloudwatch"
	@GOPATH=$(GOPATH) go get -u "gopkg.in/redis.v1"

vendor-deps: deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt dispatcher/*.go
	go fmt receiver/*.go
	go fmt *.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/wof-chatterboxd cmd/wof-chatterboxd.go
