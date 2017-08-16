CWD=$(shell pwd)
GOPATH := $(CWD)

build:	fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-flags/
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-flags/existential
	cp *.go src/github.com/whosonfirst/go-whosonfirst-flags
	cp existential/*.go src/github.com/whosonfirst/go-whosonfirst-flags/existential/
	# cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   rmdeps

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	if test ! -d src; then mkdir src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt existential/*.go

bin:	self
