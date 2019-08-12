OUT := apiserver
VERSION := $(shell git describe --always --long --dirty)
API_PKG := github.com/tarkov-database/rest-api/model/api

all: run

bin:
	go build -v -o ${OUT} -ldflags="-X ${API_PKG}.Version=${VERSION}"

lint:
	revive -config revive.toml -formatter stylish ./...

fmt:
	go fmt ./...

run: bin
	./${OUT}

clean:
	-@rm ${OUT} ${OUT}-v*
