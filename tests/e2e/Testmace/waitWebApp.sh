#!/bin/bash

HTTP_STATUS="$(curl -IL --silent http://localhost:8080 | grep HTTP | grep -P '\d\d\d' -o )"; 
echo "$HTTP_STATUS"
while [ "$HTTP_STATUS" != "404" ]; do
    echo "$HTTP_STATUS"
    echo "service can't return 404, waiting..."
    sleep 1
    HTTP_STATUS="$(curl -IL --silent http://localhost:8080 | grep HTTP | grep -P '\d\d\d' -o )"; 
done
echo "DONE $HTTP_STATUS"