GOFLAGS  := -v -tags='osusergo,netgo,static'
LDFLAGS	 := -ldflags='-s -w'

.PHONY: build-cross dist build mod clean run help docker

name		    := go-markdown-to-notion
linux_name	:= $(name)-linux-amd64
darwin_name	:= $(name)-darwin-amd64
go_version := $(shell cat $(realpath .go-version))

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

build: ## go build
	go build -o bin/$(name) $(LDFLAGS) *.go

test: ## go test
	go test -v $$(go list ./... | grep -v /vendor/)

clean: ## remove bin/*
	rm -f bin/*
