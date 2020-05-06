#!/bin/bash

set -euxo pipefail

echo "Starting to run dockerd entrypoint"
dockerd-entrypoint.sh &


for i in {1..5}
do
    echo "checking if docker started running for the agent."
    if docker version
    then
        echo "docker is now running"
        xene agent
    fi
    sleep 3
done

echo "Docker has not started yet, exitting"
exit 1
