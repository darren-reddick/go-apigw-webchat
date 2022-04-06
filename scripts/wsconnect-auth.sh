#!/usr/bin/env bash -e

source "$(dirname $0)/funcs.sh"

script="${0##*/}"
usage="USAGE: $script [sls project stage]

Example: $script dev"

# syntax checks
[ $# -ne 1 ] && printf "ERROR: Not enough arguments passed.\n\n$usage\n" && exit 1
[ "${AWS_PROFILE}" == "" ] && printf "ERROR: AWS_PROFILE environment variable required to be set for authentication.\n" && exit 1

# set the serverless stage from command line
stage=$1

required() {
    while [ $# -ne 0 ]
    do
        [ $(which $1 2>/dev/null) ] || { echo "$1 is required and must be installed"; exit 1; }
        shift
    done
}

# check required bins are installed
required jq node sls

# use serverless manifest to fetch the ServiceEndpointWebsocket name
wssurl=$(node ./scripts/get-websocket-url.js --stage ${stage})

echo "Building signed url for $wssurl using aws profile [${AWS_PROFILE}]"

# Build path to sign script and ensure it exists
SIGN_SCRIPT="$(dirname $0)/aws4-sign.js"
[ ! -f ${SIGN_SCRIPT} ] && { echo "ERROR: Script $SIGN_SCRIPT does not exist"; exit 1;  }

# connect using wscat
npx wscat -c $(node $SIGN_SCRIPT --url ${wssurl})
