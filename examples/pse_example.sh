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
        --arg label2 "${3}" \
        '{
            "value": $value | tonumber,
            "labels": [
                { "label1": $label1 },
                { "label2": $label2 }
            ]
        }'
}

for i in $(seq 1 10)
do
    add_gauge_vector ${i} "Label1 for ${i}" "Label2 for ${i}"
done