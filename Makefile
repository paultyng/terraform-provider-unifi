TEST         ?= ./...
TESTARGS     ?=
TEST_COUNT   ?= 1
TEST_TIMEOUT ?= 10m

.PHONY: default
default: build

.PHONY: build
build:
	go install

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v -count $(TEST_COUNT) -timeout $(TEST_TIMEOUT) $(TESTARGS)
