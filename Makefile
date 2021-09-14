export UNIFI_VERSION  ?= v6
export UNIFI_USERNAME ?= tfacctest
export UNIFI_EMAIL    ?= tfacctest@example.com
export UNIFI_PASSWORD ?= tfacctest1234

TEST     ?= ./...
TESTARGS ?=

.PHONY: default
default: build

.PHONY: build
build:
	go install

.PHONY: testacc
testacc:
	TF_ACC=1 UNIFI_ACC_WLAN_CONCURRENCY=4 UNIFI_API=https://localhost:8443 UNIFI_INSECURE=true go test $(TEST) -v $(TESTARGS)

.PHONY: testacc-up
testacc-up:
	docker-compose up --detach unifi

	@echo -n "Waiting for container"
	@until test -n "$$(docker ps --filter id=$$(docker-compose ps --quiet unifi) --filter health=healthy --quiet)"; do echo -n .; sleep 1; done
	@echo

	@echo "Bootstrapping Unifi controller"
	docker-compose up --abort-on-container-exit bootstrap

.PHONY: testacc-down
testacc-down:
	docker-compose down
