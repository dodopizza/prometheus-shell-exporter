#!/bin/bash

template='[
    {
        "value": 0,
        "labels": [
            { "one": "l_one" },
            { "two": "l_two" }
        ]
    }
]'

echo $template | jq
# jq -n '[1,2,3]'
