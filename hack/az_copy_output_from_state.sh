#!/usr/bin/env bash

# az_copy_output_from_state.sh
#
# script to copy the outputs from state file to another blob container
# used to share state between subscriptions without sharing the resources
# (and potentially secrets)

usage(){
    echo "sh az_copy_output_from_state.sh <subscription> <storage_account_name> <container_name> <key>"
    exit 0
}

if [ -z "$1" ]; then
    usage
fi

command -v az >/dev/null 2>&1 || { echo >&2 "The az cli is required to run this script."; exit 1; }
command -v jq >/dev/null 2>&1 || { echo >&2 "The jq command is required to run this script."; exit 1; }

SUBSCRIPTION=$1
STORAGE_ACCOUNT_NAME=$2
CONTAINER_NAME=$3
KEY=$4

terraform state pull | jq 'del(.resources)' > output.tfstate

if [ -z "$STORAGE_ACCOUNT_NAME" ]; then
    STORAGE_ACCOUNT_NAME=$(jq -r '.backend.config.storage_account_name' .terraform/terraform.tfstate)
fi

if [ -z "$CONTAINER_NAME" ]; then
    CONTAINER_NAME=$(jq -r '.backend.config.container_name' .terraform/terraform.tfstate)
fi

if [ -z "$KEY" ]; then
    KEY=$(jq -r '.backend.config.key' .terraform/terraform.tfstate)
fi

az storage blob upload -f output.tfstate --subscription $SUBSCRIPTION --account-name $STORAGE_ACCOUNT_NAME --container-name $CONTAINER_NAME --name $KEY