OUT := apiserver
VERSION := $(shell git describe --always)
REPO_PATH := tarkov-database/rest-api
API_PKG := github.com/${REPO_PATH}/model/api
IMAGE_TAG := $(shell git describe --abbrev=0 | cut -c2-)

all: run

bin:
	go build -v -o ${OUT} -ldflags="-X ${API_PKG}.Version=${VERSION}"

image:
	docker build -t docker.pkg.github.com/${REPO_PATH}/rest-api:${IMAGE_TAG} .

lint:
	revive -config revive.toml -formatter stylish ./...

fmt:
	go fmt ./...

run: bin
	./${OUT}

clean:
	-@rm ${OUT} ${OUT}-v*
