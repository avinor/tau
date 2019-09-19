# Multiple backends

Deploying with a single backend is a simple scenario that works if there is only one environment or each environment do not share any information between them. In most cases though it is required to have one remote storage per environment and maybe even share state between them. Some resources might be shared between environments and when deploying to other environments they need to have access to shared state.

## Setup

In this scenario the setup is with 3 environments:

- shared
- production
- test

Shared environment contain some networking setup, log analytics etc that the other environments need access to. However the CI / CD pipeline for test / production should only have access to it's own environment. Production should not be able to read anything from shared or test. If it can read across remote storage it could read other secrets from other state.

Challenge is how can CI / CD pipeline only have access to it's own environment while also reading output from shared?

## Solution

1. Remove all resources from shared state and create a copy that only contains output variables

2. Copy this filtered state file to other remote storages

3. For production / test access shared output variables from it's own state where only output variables are stored

## Share state

To share the state with tau run a post deploy script that copies the state from shared to the other environments. In the shared deployment files add a hook that runs on finish after running apply command.

`shared.hcl`

```terraform
hook "upload_output_state_to_prod" {
    trigger_on = "finish:apply"
    script = "https://raw.githubusercontent.com/avinor/tau/master/hack/az_copy_output_from_state.sh"
    args = ["subscription-id", "storage-account-name"]
    working_dir = module.path
}
```

This script runs in the module folder and starts by running

```bash
terraform state pull | jq 'del(.resources)' > output.tfstate
```

This simply pulls the state and remove all resources from it so only outputs remain. Once that is done it will upload state to storage account (and subscription) defined in arguments.

Script also supports 2 additional arguments to set container_name and key where to store output state. If those arguments are not set it will use same container_name and key as script just running.

In addition to configuration above it would be necessary to define the `backend` block and possibly a prepare hook.

`shared.hcl`

```terraform
backend "azurerm" {
    storage_account_name = "tfstatedevsa"
    container_name       = "state"
    key                  = "storage-account.tfstate"
}

hook "set_access_key" {
    trigger_on = "prepare"
    script = "https://raw.githubusercontent.com/avinor/terraform-azurerm-remote-backend/master/set-access-keys.sh"
    args = ["tfstatedev"]
    set_env = true
}
```

This will now use remote storage account `tfstatedevsa` and retrieve access key to this account before executing.

## Use shared state

To use the shared state in another deployment, that uses another backend, it can now read directly from its own backend.

`production.hcl`

```terraform
dependency "shared" {
    source = "../shared/shared.hcl"
    backend "azurerm" {
        storage_account_name = "storage-account-name"
    }
}

backend "azurerm" {
    storage_account_name = "storage-account-name"
    container_name       = "state"
    key                  = "prod.tfstate"
}

hook "set_access_key" {
    trigger_on = "prepare"
    script = "https://raw.githubusercontent.com/avinor/terraform-azurerm-remote-backend/master/set-access-keys.sh"
    args = ["storage-account-name"]
    set_env = true
}
```

This does several things. `backend` block defines a new backend to write remote state to, and hook block creates an access key for this block. The dependency block reference the shared deployment file, but overrides the storage_account_name to be the same as current backend. It does not override the key or container_name as those are by default set to same. Doing this it will now read the state for dependency from same remote storage as its currently using, but that state only contains the output variables and nothing else.

## Backend

Follow guide in [single backend](./single_backend.md) file on how to create a single backend. For each environment create a new backend.