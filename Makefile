OUT := perch-interactive-challenge-microservice
PKG := ./cmd
DOCKERFILE := ./build/package/Dockerfile
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: run

docker-build:
	docker build -f ${DOCKERFILE} -t ${OUT} .

deps:
	GO111MODULES=on go get -v ${PKG}

build:
	GO111MODULES=on go build -v -o ${OUT} ${PKG}

test:
	GO111MODULES=on go test -short ${PKG_LIST}

vet:
	GO111MODULES=on go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

out:
	@echo ${OUT}-v${VERSION}

protos:
	protoc -I ./protos ./protos/event.proto --go_out=./core/protos

.PHONY: run protos build docker-build vet lint out deps