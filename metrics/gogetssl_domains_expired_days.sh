#!/bin/bash
set -uo pipefail

V1_URL_PREFIX=https://my.gogetssl.com/api/
V1_API_USER=
V1_API_PASS=

# Usage: serverErrorHandler <Server responce>
function serverErrorHandler {
    [ $# -ne 1 ] && return 99
    local result=${1}
    local error=$(echo "${result}" | jq -r '.error')
    local errorMessage=$(echo "${result}" | jq -r '.message')
    local errorDescription=$(echo "${result}" | jq -r '.description')
    local resultCode=0
    if [[ ${error} == true ]]; then
        case "${message}" in
        "Wrong username/password" ) resultCode=101
        ;;
        "Auth key required for this method" ) resultCode=102
        ;;
        * ) resultCode=100
        ;;
        esac
    fi
    [[ ${error} == null && ${result} != null && ${result} != "" ]] && echo "${result}"
    return ${resultCode}
}

# Usage: getApiToken <URL prefix> <API user> <API password>
function getApiToken {
    [ $# -ne 3 ] && return 99
    local urlPrefix=${1}
    local apiUser=${2}
    local apiPassword=${3}
    local resultCode=0
    local result=$(serverErrorHandler $(
        curl --silent --location -g --request POST "${urlPrefix}/auth" \
            --header 'Content-Type: application/x-www-form-urlencoded' \
            --data-urlencode "user=${apiUser}" \
            --data-urlencode "pass=${apiPassword}"
        )
    )
    resultCode=$?
    [ ${resultCode} -eq 0 ] && (echo "${result}" | jq -r '.key')
    return ${resultCode}
}

# Usage: getActiveOrders <URL prefix> <API token>
function getActiveOrders {
    [ $# -ne 2 ] && return 99
    local urlPrefix=${1}
    local apiToken=${2}

    local result=$(serverErrorHandler $(curl --silent --location -g --request GET "${urlPrefix}/orders/ssl/all?auth_key=${apiToken}&limit=999"))
    resultCode=$?
    [ ${resultCode} -eq 0 ] && (echo "${result}" | jq -r '.orders[] | select(.status == "active").order_id')
    return ${resultCode}
}

# Usage: ordersToJson <active orders> <URL prefix> <API token>
function ordersToJson {
    [ $# -ne 3 ] && return 99
    local activeOrders=${1}
    local urlPrefix=${2}
    local apiToken=${3}
    (
        for orderId in ${activeOrders}; do
            local result=$(curl --silent --location -g --request GET "${urlPrefix}/orders/status/${orderId}?auth_key=${apiToken}")
            local error=$(echo "${result}" | jq -r '.error')
            [[ ${error} != true ]] && (echo -n "${result}" |
                jq '{ 
                    value: (((((.valid_till + "T00:00:00Z")|fromdate) - now) / 60 / 60 / 24)|round), 
                    labels: { 
                        domain: (.domain, .san[].san_name) 
                    } 
                }'
            )
        done
    ) | jq -s
    return 0
}

function main {
    local result=""
    local resultCode=0
    case "${1:-}" in
    "token" )
        result=$(getApiToken ${V1_URL_PREFIX} ${V1_API_USER} ${V1_API_PASS})
        resultCode=$?
    ;;
    * )
        # Get token if not set
        if [[ ${V1_API_TOKEN:-} == "" ]]; then 
            echo "--1--"
            export V1_API_TOKEN=$(main "token")
            resultCode=$?
        fi
        # Get list of orders
        if [ ${resultCode} -eq 0 ]; then
            echo "--2--"
            activeOrders=$(getActiveOrders ${V1_URL_PREFIX} ${V1_API_TOKEN}) 
            resultCode=$?
            # If get server error, trying to get token
            if [ ${resultCode} -ne 0 ]; then
                echo "--3--"
                export V1_API_TOKEN=$(main "token")
                resultCode=$?
                # If success token query
                if [ ${resultCode} -eq 0 ]; then
                    echo "--4--"
                    activeOrders=$(getActiveOrders ${V1_URL_PREFIX} ${V1_API_TOKEN}) 
                    resultCode=$?
                fi
            fi
        fi
        # If order list is present
        if [ ${resultCode} -eq 0 ]; then
            echo "--5--"
            result=$(ordersToJson "${activeOrders[*]}" "${V1_URL_PREFIX}" "${V1_API_TOKEN}")
            resultCode=$?
        fi
    ;;
    esac

    [ ${resultCode} -eq 0 ] && echo ${result} && return 0
    [ ${resultCode} -ne 0 ] && return ${resultCode}
}

main $@
