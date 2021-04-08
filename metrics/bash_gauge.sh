#!/bin/bash

metric_result=""

trap on_exit EXIT

function on_exit(){
    echo "${metric_result}" | jq -c -s '.'
}

# Usage: add_gauge_vector
function add_gauge_vector(){
    metric_result+=$(gauge_vector "$@")
}

# Usage: gauge_vector 0 "One" "Two"
function gauge_vector(){
    jq -n \
        --arg value "${1}" \
        --arg label1 "${2}" \
        '{
            "value": $value | tonumber,
            "labels": {
                "example": $label1,
            }
        }'
}

for i in $(seq 1 10)
do
    add_gauge_vector ${RANDOM} "for_${i}"
done