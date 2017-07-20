CWD=$(shell pwd)
GOPATH := $(CWD)

build:	fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson-v2/feature
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson-v2/properties
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson-v2/utils
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson-v2/whosonfirst
	cp -r geojson/*.go src/github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson/
	cp -r feature/*.go src/github.com/whosonfirst/go-whosonfirst-geojson-v2/feature/
	cp -r properties/ src/github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/
	cp -r utils/*.go src/github.com/whosonfirst/go-whosonfirst-geojson-v2/utils/
	cp -r whosonfirst/*.go src/github.com/whosonfirst/go-whosonfirst-geojson-v2/whosonfirst/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/skelterjohn/geom"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt geojson/*.go
	go fmt feature/*.go
	go fmt properties/*/*.go
	go fmt utils/*.go
	go fmt whosonfirst/*.go

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-dump cmd/wof-geojson-dump.go
