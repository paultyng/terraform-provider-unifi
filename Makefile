TEST       ?= ./...
TESTARGS   ?=
TEST_COUNT ?= 1

.PHONY: default
default: build

.PHONY: build
build:
	go install

.PHONY: testacc
testacc:
	TF_ACC=1 UNIFI_ACC_WLAN_CONCURRENCY=3 UNIFI_API=https://localhost:8443 UNIFI_INSECURE=true UNIFI_USERNAME=admin UNIFI_PASSWORD=admin go test $(TEST) -v -count=$(TEST_COUNT) $(TESTARGS)

.PHONY: testacc-up
testacc-up:
	docker compose up --detach --wait

.PHONY: testacc-down
testacc-down:
	docker compose down
