#! /bin/bash

ver="5.12.35"
keys=$(jq -r keys[] "$ver/Setting.json")

while IFS= read -r key; do
    fn="$(echo $key | sed -r 's/(^|_)([a-z])/\U\2/g')"
    echo "... $key $fn ..."
    jq ".$key" "$ver/Setting.json" >> "$ver/Setting$fn.json"
done <<< "$keys"