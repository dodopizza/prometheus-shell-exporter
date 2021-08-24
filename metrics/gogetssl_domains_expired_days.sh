#!/bin/bash
V1_URL_PREFIX=https://my.gogetssl.com/api/
V1_API_USER=
V1_API_PASS=
V1_API_TOKEN=

# Input parameters:
# $1 - URL prefix
# $2 - API user
# $3 - API password
# Output:
# 0 - Everything is correct
# 1 - invalid arguments number
# 2 - empty answer from server
function getApiToken {
    if [ $# -ne 3 ]; then
        return 1
    else
        result=$(
            curl --silent --location -g --request POST "${1}/auth" \
                --header 'Content-Type: application/x-www-form-urlencoded' \
                --data-urlencode "user=${2}" \
                --data-urlencode "pass=${3}" |
                jq -r '.key' 2> /dev/null
        )
        if [[ ${result} == "" ]]; then
            return 2
        else
            echo "${result}"
            return 0
        fi
    fi
}

# Input parameters:
# $1 - URL prefix
# $2 - API token
# Output:
# 0 - Everything is correct
# 1 - invalid arguments number
# 2 - empty answer from server
function getActiveOrders {
    if [ $# -ne 2 ]; then
        return 1
    else
        result=$(
            curl --silent --location -g --request GET "${1}/orders/ssl/all?auth_key=${2}&limit=999" |
                jq -r '.orders[] | select(.status == "active").order_id' 2> /dev/null
            )
        if [[ ${result} == "" ]]; then
            return 2
        else
            echo "${result}"
            return 0
        fi
    fi
}

# Input parameters:
# $1 - active orders
# $2 - URL prefix
# $3 - API token
# Output:
# 0 - Everything is correct
# 1 - invalid arguments number
# 2 - empty answer from server
function ordersToJson {
    if [ $# -ne 3 ]; then
        return 1
    else
        for order_id in $1; do
            result+=$(curl --silent --location -g --request GET "${2}/orders/status/${order_id}?auth_key=${3}" | 
                jq '{ 
                    value: (((((.valid_till + "T00:00:00Z")|fromdate) - now) / 60 / 60 / 24)|round), 
                    labels: { 
                        domain: (.domain, .san[].san_name) 
                    } 
                }')
        done
        if [[ ${result} == "" ]]; then
            return 2
        else
            echo ${result}
            return 0
        fi
    fi
}

if [[ $1 == "token" ]]; then
    echo $(getApiToken ${V1_URL_PREFIX} ${V1_API_USER} ${V1_API_PASS})
else
    if [[ ${V1_API_TOKEN} != "" ]]; then
        activeOrders=$(getActiveOrders ${V1_URL_PREFIX} ${V1_API_TOKEN})
        activeOrdersCode=$?
    fi

    if [[ ${activeOrdersCode} == 2 || ${V1_API_TOKEN} == "" ]]; then
        token=$(getApiToken ${V1_URL_PREFIX} ${V1_API_USER} ${V1_API_PASS})
        tokenCode=$?
        if [[ ${tokenCode} == 0 ]]; then
            V1_API_TOKEN=${token}
            activeOrders=$(getActiveOrders ${V1_URL_PREFIX} ${V1_API_TOKEN})
            activeOrdersCode=$?
        fi
    fi
    echo $(ordersToJson "${activeOrders[*]}" "${V1_URL_PREFIX}" "${V1_API_TOKEN}")
fi
