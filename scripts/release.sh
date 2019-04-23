#!/bin/bash

body='{
"request": {
"branch":"master"
}}'

curl -vvv -s -X POST \
 -H "Content-Type: application/json" \
 -H "Accept: application/json" \
 -H "Travis-API-Version: 3" \
 -H "Authorization: token $(travis token --org)" \
 -d "$body" \
 https://api.travis-ci.org/repo/PremiereGlobal%2Fstim/requests
