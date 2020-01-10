#! /bin/bash

set -eou pipefail

# Local Administrator
# Username: tfacctest
# Password: tfacctest1234
# Email: tfacctest@example.com

if test $# -eq 0; then
  echo "please specify either 'start' or 'test'"
  exit 1
fi

case "$1" in
  "start")
    docker run --rm --init -d \
      -p 8080:8080 \
      -p 8443:8443 \
      -p 3478:3478/udp \
      -p 10001:10001/udp \
      -e TZ='America/New_York' \
      -v $(pwd)/testdata/unifi:/unifi \
      --name unifi \
      jacobalberty/unifi:stable

    echo "Waiting for login page..."
    timeout 300 bash -c 'while [[ "$(curl --insecure -s -o /dev/null -w "%{http_code}" https://localhost:8443/manage/account/login)" != "200" ]]; do sleep 5; done'
    echo "Controller running."
    ;;
  "test")
    TF_ACC=1 \
    UNIFI_USERNAME=tfacctest \
    UNIFI_PASSWORD=tfacctest1234 \
    UNIFI_API="https://localhost:8443/api/" \
    UNIFI_ACC_WLAN_CONCURRENCY="4" \
    go test -v -cover ./internal/provider
    ;;
  *)
    echo "unrecognized command"
    exit 1
    ;;
esac