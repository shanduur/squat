#!/bin/env bash

FIRST_RUN=/etc/squat/first

if [[ -f $FIRST_RUN ]]; then
    /bin/gob-generator -i /etc/squat/data/data.json -o /etc/squat/data/data.gob
    touch $FIRST_RUN
fi

/bin/squat