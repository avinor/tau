[![CircleCI](https://circleci.com/gh/avinor/tau.svg?style=svg)](https://circleci.com/gh/avinor/tau)
[![Go Report Card](https://goreportcard.com/badge/github.com/avinor/tau)](https://goreportcard.com/report/github.com/avinor/tau)
[![GoDoc](https://godoc.org/github.com/avinor/tau?status.svg)](https://godoc.org/github.com/avinor/tau)

# Tau

Tau (Terraform Avinor Utility) is a thin wrapper for [Terraform](https://www.terraform.io/) that makes terraform execution and secret handling easier. It tries to conform to the principal of keeping code DRY (Don't Repeat Yourself) and bases all executions on terraform modules. Since terraform already provides a lot of excellent features it tries to not change that too much but use those features best possible.

There are also other good tools available, see the [comparison](#comparison) at the end.

**NOTE: This is still in development.**

## Installation

1. Tau requires [terraform](https://www.terraform.io/) 0.12+, download and install first
2. Download tau from [Release page](https://github.com/avinor/tau/releases) for your OS
3. Rename file to `tau` and add it to your `PATH`

## Highlights

**DRY:** All executions are based on using terraform modules. Tau configuration file just describes how to execute / deploy the modules.

**Dependency handling:** Tau handles dependencies between different modules. Executing them in correct order and passing output from one module as input to another.

**Secret handling:** Recommended way to deal with secrets is passing them in as input variables. Using the build-in data sources in terraform tau can send anything as input variables to terraform.

**Backend:** Backend configuration is taken out of modules and defined in tau configuration.

## Configuration

Any files named `.hcl` or `.tau` are read, where each file is one deployment of module.

```terraform
// One or many hooks that can trigger on prepare or finish
hook "set_access_key" {
    trigger_on = "prepare"
    command = "./set_access_key.sh"
    set_env = true
}

// One or more dependencies
dependency "vnet" {
    source = "./vnet.hcl"
}

dependency "logs" {
    source = "./logs.hcl"
}

// One or more data blocks, support any terraform data block
data "azurerm_key_vault_secret" "secret" {
    name = "my-secret"
    key_vault_id = "/subscriptions/xxxx-xxxx-xxxx-xxxx/resourceGroups/secrets-rg/providers/Microsoft.KeyVault/vaults/terraform-secrets-kv"
}

data "azurerm_key_vault_secret" "sp" {
    name = "sp-secret"
    key_vault_id = "/subscriptions/xxxx-xxxx-xxxx-xxxx/resourceGroups/secrets-rg/providers/Microsoft.KeyVault/vaults/terraform-secrets-kv"
}

// Set environment variables for all terraform commands
environment_variables {
    ARM_SUBSCRIPTION_ID = "xxxx-xxxx-xxxx-xxxx"
}

backend "azurerm" {
    storage_account_name = "terraformstatesa"
    container_name       = "state"
    key                  = "westeurope/${source.name}.tfstate"
}

// Define which module to deploy.
module {
    source = "avinor/kubernetes/azurerm"
    version = "1.0.0"
}

// Input variables to module
inputs {
    name = "example"
    resource_group_name = "example-rg"
    service_cidr = "10.241.0.0/24"
    kubernetes_version = "1.13.5"
    log_analytics_workspace_id = dependency.logs.outputs.resource_id

    service_principal = {
        client_id = data.azurerm_key_vault_secret.secret.value
        client_secret = data.azurerm_key_vault_secret.secret.value
    }

    azure_active_directory = {
        client_app_id = data.azurerm_key_vault_secret.sp.value
        server_app_id = data.azurerm_key_vault_secret.sp.value
        server_app_secret = data.azurerm_key_vault_secret.sp.value
    }

    agent_pools = [
        {
            name = "ipt"
            vm_size = "Standard_D2_v3"
            vnet_subnet_id = dependency.vnet.outputs.subnets.aks
        },
    ]
}
```

### hook

One or more hooks that triggers on specific events during deployment. It can read the output from command run and set environment variables for terraform, for instance access keys etc.

To read output and set environment variables set `set_env` = true. It will read all output in format "key = value" and add them to the environment when running terraform.

If `fail_on_error` is set it will accept any failures from command and continue executing terraform commands. Default value is false and it will stop all executions.

To optimize execution and not run same command multiple times (for instance retrieving same access key) it caches output from every command and reuses cached value if called multiple times in same run. To disable cache set `disable_cache` = true.

attribute | Description
----------|------------
trigger_on    | Event to trigger hook on. Possible values are "prepare" and "finish"
command       | Command to execute, should not include arguments
args          | Arguments to send to command
set_env       | If true it will read output in format "key = value"
fail_on_error | Fail on error or continue running ignoring error
disable_cache | Disable cache and make sure command is run every time

## Comparison

A short comparison to other similar tools and why we decided to create and use `tau`.

### [Terragrunt](https://github.com/gruntwork-io/terragrunt)

Terragrunt is a nice tool to reuse modules for multiple deployments, however we disagreed on some of the design choices and also wanted to have a better way to handle secrets.

We started using terragrunt but with the 0.12 release of terraform it stopped working until terragrunt updated to latest terraform syntax. This was a weakness in terragrunt being too dependant on terraform syntax. It now uses a newer syntax and therefore does not have this problem anymore.

#### Backend defined in module

Terragrunt requires that the backend is defined in the module because it uses `terraform init -from-module ...`, so backend block has to exist to override.

Tau instead uses `go-getter` directly to download the module and then creates an override file defining the backend. It will run `terraform init` afterwards when all overrides have been defined.

#### Dependencies

Terragrunt recommends using a `terraform_remote_state` data source to retrieve dependencies in modules. We disagree that this should be in the module and rather a part of input. See documentation for more details how this is solved in tau.

#### Secrets

Terragrunt don't have any way to handle secrets. It supports sending environment variables to commands, but cannot retrieve secrets from an Azure Key Vault for instance.

### [Astro](https://github.com/uber/astro/)

Astro is a tool created by Uber to manage terraform executions.

#### YAML

There is nothing wrong with yaml format, but terraform has decided to use its own format called hcl. Since tool is working on terraform code we find it more natural to also use hcl instead of yaml.

#### One big file

Astro defines all modules in one big file instead of supporting separate files for each deployment.

#### Secrets

No secret handling, does not support retrieving values from a secret store.
