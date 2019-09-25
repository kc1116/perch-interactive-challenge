OUT := perch-iot-pubsub
PKG := .
DOCKERFILE := ./build/Dockerfile
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: run

docker-build:
	docker build -f ${DOCKERFILE} -t ${OUT} .

run_network:
	docker-compose -f ./build/docker-compose.yaml up --abort-on-container-exit

deps:
	GO111MODULES=on go get ${PKG}

build:
	GO111MODULES=on go build -o ${OUT} ${PKG}

deploy: deps build
	mv ./perch-iot-pubsub /usr/bin/perch-iot-pubsub

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

.PHONY: run protos build docker-build vet lint out deps deploy run_network