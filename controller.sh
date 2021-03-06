#! /bin/bash

set -eou pipefail

default_tag="stable-6"

# Local Administrator
# Username: tfacctest
# Password: tfacctest1234
# Email: tfacctest@example.com
DOCKER_HTTP_PORT="${DOCKER_HTTP_PORT:-8080}"
DOCKER_HTTPS_PORT="${DOCKER_HTTPS_PORT:-8443}"
DOCKER_STUN_PORT="${DOCKER_STUN_PORT:-3478}"
DOCKER_AIRCONTROL_PORT="${DOCKER_AIRCONTROL_PORT:-10001}"

if test $# -eq 0; then
  echo "please specify either 'start' or 'test'"
  exit 1
fi

case "$1" in
  "start")
    docker run --rm --init -d \
      -p ${DOCKER_HTTP_PORT}:8080 \
      -p ${DOCKER_HTTPS_PORT}:8443 \
      -p ${DOCKER_STUN_PORT}:3478/udp \
      -p ${DOCKER_AIRCONTROL_PORT}:10001/udp \
      -e TZ='America/New_York' \
      --name unifi \
      jacobalberty/unifi:${2:-$default_tag}

    echo "Waiting for login page..."
    timeout 300 bash -c 'while [[ "$(curl --insecure --location -s -o /dev/null -w "%{http_code}" '"https://localhost:${DOCKER_HTTPS_PORT}/manage/account/login"')" != "200" ]]; do sleep 5; done'
    echo "Controller running."

    echo "Bootstrapping controller..."
    function unifi_curl() {
      curl --data "${2}" --header 'Content-Type application/json' --insecure "https://localhost:${DOCKER_HTTPS_PORT}${1}"
    }
    unifi_curl /api/cmd/sitemgr '{"cmd":"add-default-admin","name":"tfacctest","email":"tfacctest@example.com","x_password":"tfacctest1234"}'
    unifi_curl /api/set/setting/super_identity '{"name":"UniFi Network"}'
    unifi_curl /api/cmd/system '{"cmd":"set-installed"}'
    ;;
  "test")
    TF_ACC=1 \
    UNIFI_USERNAME=tfacctest \
    UNIFI_PASSWORD=tfacctest1234 \
    UNIFI_API="https://localhost:${DOCKER_HTTPS_PORT}/" \
    UNIFI_ACC_WLAN_CONCURRENCY="4" \
    UNIFI_INSECURE="true" \
    go test -v -cover -count 1 ./internal/provider
    ;;
  "stop")
    docker stop unifi
    ;;
  "update")
    docker pull jacobalberty/unifi:${2:-$default_tag}
    ;;
  *)
    echo "unrecognized command"
    exit 1
    ;;
esac
