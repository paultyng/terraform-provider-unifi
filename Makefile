TEST       ?= ./...
TESTARGS   ?=
TEST_COUNT ?= 3

.PHONY: default
default: build

.PHONY: build
build:
	go install

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v -count=$(TEST_COUNT) $(TESTARGS) -timeout 20m
