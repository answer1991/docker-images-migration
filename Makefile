SHELL          = /bin/bash

BASE_IMAGE     = golang:1.7.5-alpine

SRC = github.com/answer1991/docker-images-migration
VERSION = $(shell cat VERSION)
GITCOMMIT = $(shell git log -1 --pretty=format:%h)
BUILD_TIME = $(shell date "+%Y-%m-%d_%H-%M-%S")

default: build-linux

build-mac:
	docker run --rm -v $(shell pwd):/go/src/${SRC} -w /go/src/${SRC} -e CGO_ENABLED=0 -e GOOS=darwin ${BASE_IMAGE}  go build -a -v -ldflags "-X ${SRC}/cmd.Platform=MacOs -X ${SRC}/cmd.Version=${VERSION} -X ${SRC}/cmd.GitCommit=${GITCOMMIT} -X ${SRC}/cmd.BuildTime=${BUILD_TIME}"

build-linux:
	docker run --rm -v $(shell pwd):/go/src/${SRC} -w /go/src/${SRC} -e CGO_ENABLED=0 ${BASE_IMAGE}  go build -a -v -ldflags "-X ${SRC}/cmd.Platform=Linux -X ${SRC}/cmd.Version=${VERSION} -X ${SRC}/cmd.GitCommit=${GITCOMMIT} -X ${SRC}/cmd.BuildTime=${BUILD_TIME}"