## Stack Name
STACK_NAME := backlog-to-slack-dm

## Install library for production
deps:
	go get -u ./...
.PHONY: deps

## Install library for development
devel-deps: deps
	GO111MODULE=off go get \
		golang.org/x/lint/golint \
		honnef.co/go/tools/staticcheck \
		github.com/kisielk/errcheck \
		golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow \
		github.com/securego/gosec/cmd/gosec \
		github.com/motemen/gobump/cmd/gobump \
		github.com/Songmu/make2help/cmd/make2help
.PHONY: devel-deps

## Execute unit test
test: deps
	go test -v -count=1 -cover ./...
.PHONY: test

## Output coverage of testing
cov:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out
.PHONY: cov

## Clean up artifact
clean:
	rm -rf ./artifact/*
.PHONY: clean

## Lint
lint: devel-deps
	go vet ./...
	staticcheck ./...
	errcheck ./...
	gosec -quiet ./...
	golint -set_exit_status ./...
.PHONY: lint

## Build
build: build-backlogtoslackdm build-authorizer
.PHONY: build

## Build backlog-to-slack-dm
build-backlogtoslackdm:
	GOOS=linux GOARCH=amd64 go build -o artifact/backlog-to-slack-dm ./handlers/backlog-to-slack-dm
.PHONY: build-backlogtoslackdm

## Build authorizer
build-authorizer:
	GOOS=linux GOARCH=amd64 go build -o artifact/authorizer ./handlers/authorizer
.PHONY: build-authorizer

## SAM Validate
validate:
	sam validate
.PHONY: validate

## Deploy by SAM
deploy: build
	sam deploy
.PHONY: deploy

## Delete CloudFormation Stack
delete:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
.PHONY: delete

## Help
help:
	@make2help --all
.PHONY: help
