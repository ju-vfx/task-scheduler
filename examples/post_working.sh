#! /bin/bash

curl --request POST \
    --header "Content-Type: application/json" \
    --data @./working_job.json \
    http://localhost:8080/api/jobs