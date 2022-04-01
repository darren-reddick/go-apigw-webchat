# utility to get the websocket url
getWsUrl() {
    wssurl=$(npx sls manifest --json --stage ${1} | jq -r '.'${1}'.outputs[] | select(.OutputKey == "ServiceEndpointWebsocket").OutputValue')
    echo $wssurl
}