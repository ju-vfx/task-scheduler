#! /bin/bash

curl --request POST \
    --header "Content-Type: application/json" \
    --data @./error_job.json \
    http://localhost:8080/api/jobs