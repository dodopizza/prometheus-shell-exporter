#!/bin/bash
set -euo pipefail
V1_API_USER=
V1_API_PASS=

# Usage: getApiToken <URL prefix> <API user> <API password>
function getApiToken {
    [ $# -ne 3 ] && return 99
    local url_prefix=${1}
    local api_user=${2}
    local api_password=${3}

    curl --silent --location -g --request POST "${url_prefix}/auth" \
        --header 'Content-Type: application/x-www-form-urlencoded' \
        --data-urlencode "user=${api_user}" \
        --data-urlencode "pass=${api_password}" |
        jq -r '.key'
}

# Usage: getActiveOrders <URL prefix> <API token>
function getActiveOrders {
    [ $# -ne 2 ] && return 99
    local url_prefix=${1}
    local api_token=${2}
    curl --silent --location -g --request GET "${url_prefix}/orders/ssl/all?auth_key=${api_token}&limit=999" |
        jq -r '.orders[] | select(.status == "active").order_id' 2>/dev/null
}

# Usage: ordersToJson <active orders> <URL prefix> <API token>
function ordersToJson {
    [ $# -ne 2 ] && return 99
    local active_orders=${1}
    local url_prefix=${2}
    local api_token=${3}
    (
        for order_id in ${active_orders}; do
            curl --silent --location -g --request GET "${url_prefix}/orders/status/${order_id}?auth_key=${api_token}" |
                jq '{ 
                    value: (((((.valid_till + "T00:00:00Z")|fromdate) - now) / 60 / 60 / 24)|round), 
                    labels: { 
                        domain: (.domain, .san[].san_name) 
                    } 
                }'
        done
    ) | jq -s
}

# Usage: main
function main {
    local token="${V1_API_TOKEN:-}"
    local url_prefix="${V1_URL_PREFIX:-https://my.gogetssl.com/api/}"

    if [ -z "${token}" ]; then
        token=$(getApiToken "${url_prefix}" "${V1_API_USER}" "${V1_API_PASS}" || true)
    fi

    echo "${token}"

    [ -z "${token}" ] && (
        echo "Can't get token"
        return 1
    )

    activeOrders=$(getActiveOrders ${url_prefix} ${token})

    return 0
}

main $@

# if [[ $1 == "token" ]]; then
#     echo $(getApiToken ${V1_URL_PREFIX} ${V1_API_USER} ${V1_API_PASS})
# else
#     if [[ ${V1_API_TOKEN} != "" ]]; then
#         activeOrders=$(getActiveOrders ${V1_URL_PREFIX} ${V1_API_TOKEN})
#         activeOrdersCode=$?
#     fi

#     if [[ ${activeOrdersCode} == 2 || ${V1_API_TOKEN} == "" ]]; then
#         token=$(getApiToken ${V1_URL_PREFIX} ${V1_API_USER} ${V1_API_PASS})
#         tokenCode=$?
#         if [[ ${tokenCode} == 0 ]]; then
#             V1_API_TOKEN=${token}
#             activeOrders=$(getActiveOrders ${V1_URL_PREFIX} ${V1_API_TOKEN})
#             activeOrdersCode=$?
#         fi
#     fi
#     echo $(ordersToJson "${activeOrders[*]}" "${V1_URL_PREFIX}" "${V1_API_TOKEN}")
# fi
