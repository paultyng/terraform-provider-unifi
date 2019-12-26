#! /bin/bash

# docker run \
#   --name=unifi-controller \
#   -e PUID=1000 \
#   -e PGID=1000 \
#   -e MEM_LIMIT=1024M `#optional` \
#   -p 3478:3478/udp \
#   -p 10001:10001/udp \
#   -p 8080:8080 \
#   -p 8081:8081 \
#   -p 8443:8443 \
#   -p 8843:8843 \
#   -p 8880:8880 \
#   -p 6789:6789 \
#   -v $(pwd)/testdata/config:/config \
#   --rm \
#   linuxserver/unifi-controller:LTS

docker run --rm --init \
  -p 8080:8080 \
  -p 8443:8443 \
  -p 3478:3478/udp \
  -p 10001:10001/udp \
  -e TZ='America/New_York' \
  -v ~/testdata/unifi:/unifi \
  --name unifi \
  jacobalberty/unifi:stable
