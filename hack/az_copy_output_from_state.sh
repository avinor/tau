#!/usr/bin/env bash

# az_copy_output_from_state.sh
#
# script to copy the outputs from state file to another blob container
# used to share state between subscriptions without sharing the resources
# (and potentially secrets)

usage(){
    echo "sh az_copy_output_from_state.sh <container-name> <blob-name>"
    exit 0
}

if [ -z "$1" ]; then
    usage
fi

if [ -z "$2" ]; then
    usage
fi

command -v az >/dev/null 2>&1 || { echo >&2 "The az cli is required to run this script."; exit 1; }
command -v jq >/dev/null 2>&1 || { echo >&2 "The jq command is required to run this script."; exit 1; }

CONTAINER_NAME=$1
BLOB_NAME=$2
