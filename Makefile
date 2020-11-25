PKG_TARGET=linux
PKG_BIN=./bin/enve
PKG_TAG=$(shell git tag -l --contains HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%m:%S')

export GO111MODULE := on
# enable consistent Go 1.12/1.13 GOPROXY behavior.
export GOPROXY = https://proxy.golang.org,https://gocenter.io,direct

#######################################
############# Development #############
#######################################

install:
	@go version
	@go get -v golang.org/x/lint/golint
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	@curl -L https://git.io/misspell | sh
	@go mod download
	@go mod tidy
.ONESHELL: install

watch:
	@refresh run
.ONESHELL: watch

dev.release:
	set -e
	set -u

	@goreleaser release --skip-publish --rm-dist --snapshot
.ONESHELL: dev.release


#######################################
########### Utility tasks #############
#######################################

test:
	@go version
	@golint -set_exit_status ./...
	@go vet ./...
	@go test $$(go list ./...) \
		-v -timeout 30s -race -coverprofile=coverage.txt -covermode=atomic
.PHONY: test

coverage:
	@bash -c "bash <(curl -s https://codecov.io/bash)"
.PHONY: coverage

tidy:
	@go mod tidy
.PHONY: tidy

fmt:
	@find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

lint:
	@go version
	@./bin/golangci-lint run --tests=false --enable-all --disable=lll --disable funlen --disable godox ./...
	@./bin/misspell -error **/*
.PHONY: lint

dev_release:
	@go version
	@goreleaser release --snapshot --rm-dist
.PHONY: dev_release

build:
	@go version
	@go build -v \
		-ldflags "-s -w \
		-X 'github.com/joseluisq/enve/cmd.versionNumber=0.0.0' \
		-X 'github.com/joseluisq/enve/cmd.buildTime=$(BUILD_TIME)'" \
		-a -o bin/enve main.go
.PHONY: build


#######################################
########## Production tasks ###########
#######################################

prod.release:
	set -e
	set -u

	@go version
	@git tag $(GIT_TAG)
	@goreleaser release --rm-dist
.ONESHELL: prod.release

prod.release.ci:
	set -e
	set -u

	@go version
	@git tag $(DRONE_TAG)
	@curl -sL https://git.io/goreleaser | bash
.ONESHELL: prod.release.ci
