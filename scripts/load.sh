#!/usr/bin/env bash

TAGS=""

for i in $(cat .dev/providers.json | jq -r ".${1}[]"); do 
    TAGS+="-tags=${i} "
done

echo $TAGS
