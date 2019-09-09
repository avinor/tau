# Single subscription

When deploying resources to a single subscription it is easy to handle the remote state. This example is based on Azure but can be used for other remote state backends as well.

In a single subscription all resources are deployed to the same subscription. Used / service account that is running script therefore have access to read all resources.

## Backend

First create the backend that will keep all the state files.

```terraform
hook "set_access_key" {
    trigger_on = "prepare"
    command = "github.com/avinor/terraform-azurerm-remote-backend/setAccessKeys.sh"
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
