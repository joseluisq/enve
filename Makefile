PKG_TARGET=linux
PKG_BIN=./bin/enve
PKG_TAG=$(shell git tag -l --contains HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%m:%S')

export GO111MODULE := on

#######################################
############# Development #############
#######################################

install:
	@go version
	@go install honnef.co/go/tools/cmd/staticcheck@2023.1.3
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
	@staticcheck ./...
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
	@git tag $(GIT_TAG) --sign -m "$(GIT_TAG)"
	@goreleaser release --clean --skip=publish --skip=validate
.ONESHELL: prod.release

prod.release.ci:
	set -e
	set -u

	@go version
	@git tag $(GIT_TAG) --sign -m "$(GIT_TAG)"
	@curl -sL https://git.io/goreleaser | bash
.ONESHELL: prod.release.ci
