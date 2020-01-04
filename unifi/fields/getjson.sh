#! /bin/bash

#! /bin/bash

deburl="https://dl.ui.com/unifi/5.12.35/unifi_sysvinit_all.deb"
wkdir="$(mktemp -d)"
deb="$wkdir\unifi.deb"

curl -o "$deb" "$deburl"

mkdir -p "$wkdir/unifi"
dpkg-deb -R "$deb" "$wkdir/unifi"

# cp "$wkdir/unifi/usr/lib/unifi/lib/ace.jar" ./

# TODO: extract the JSON field files