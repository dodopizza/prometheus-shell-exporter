#!/bin/bash
V1_URL_PREFIX=https://my.gogetssl.com/api/
V1_API_USER=***
V1_API_PASS=***

V1_API_TOKEN=$(
    curl --silent --location -g --request POST "${V1_URL_PREFIX}/auth" \
        --header 'Content-Type: application/x-www-form-urlencoded' \
        --data-urlencode "user=${V1_API_USER}" \
        --data-urlencode "pass=${V1_API_PASS}" |
        jq -r '.key'
)

active_orders=$(
    curl --silent --location -g --request GET "${V1_URL_PREFIX}/orders/ssl/all?auth_key=${V1_API_TOKEN}&limit=999" |
        jq -r '.orders[] | select(.status == "active").order_id'
)

(
    for order_id in ${active_orders}; do
        curl --silent --location -g --request GET "${V1_URL_PREFIX}/orders/status/${order_id}?auth_key=${V1_API_TOKEN}" | 
        jq '{ 
            value: (((((.valid_till + "T00:00:00Z")|fromdate) - now) / 60 / 60 / 24)|round), 
            labels: { 
                domain: (.domain, .san[].san_name) 
            } 
        }'
    done 
) | jq -s
