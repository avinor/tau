# Single backend

When deploying resources with a single backend it is easy to handle the remote state. This example is based on Azure but can be used for other remote state backends as well.

In a single backend all resources that are deployed store state in same storage. State will be split in different files, but since it is using same storage a user that have access to storage will be able to read the entire state.

It is **recommended** to not store state for different environments in same storage. Users can access state across environments then.

## Configuration

To configure the backend see section below. Assuming backend is created it can be used by including a `backend` block in configuration file.

```terraform
backend "azurerm" {
    storage_account_name = "tfstatedevsa"
    container_name       = "state"
    key                  = "storage-account.tfstate"
}
```

Including this in every module is not very DRY code, so auto imported files can be used for this. Create a file ending in `_auto.hcl` and add backend section. When creating new module deployments create files in same folder and they will all include the auto file.

Since each module deployment needs its own file the auto file has to be changed a bit. Defining it in every module the `key` property can be set for each module. In auto file it can be set to module deployment file name.

```terraform
backend "azurerm" {
    storage_account_name = "tfstatedevsa"
    container_name       = "state"
    key                  = "${source.name}.tfstate"
}
```

Now if creating a file called `storage-account.hcl` that does not have any `backend` block will include the auto file and set key to `storage-account.tfstate`.

## Hook

Using a hook can be helpful to generate access keys before running any commands. The Azure remote state repository includes a script to set access key. In auto file, or individual deployments, create a hook section to execute the script on prepare stage.

```terraform
hook "set_access_key" {
    trigger_on = "prepare"
    script = "https://raw.githubusercontent.com/avinor/terraform-azurerm-remote-backend/master/set-access-keys.sh"
    args = ["tfstatedev"]
    set_env = true
}
```

## Backend

This configuration using module for Azure remote state, but could be any other module. Follow the guide in [Azure remote state README](https://github.com/avinor/terraform-azurerm-remote-backend/blob/master/README.md) on how to create the storage.

`remote_state.hcl`

```terraform
hook "set_access_key" {
    trigger_on = "prepare"
    command = "https://raw.githubusercontent.com/avinor/terraform-azurerm-remote-backend/master/set-access-keys.sh"
    args = ["tfstatedev"]
    set_env = true
    fail_on_error = false
}

environment_variables {
    ARM_SUBSCRIPTION_ID = "00000000-0000-0000-0000-000000000000"
}

backend "azurerm" {
    storage_account_name = "tfstatedevsa"
    container_name       = "state"
    key                  = "backend/dev.tfstate"
}

module {
    source = "avinor/remote-backend/azurerm"
    version = "1.0.3"
}

inputs {
    name = "tfstatedev"
    resource_group_name = "terraform-rg"
    location = "westeurope"
}
```
