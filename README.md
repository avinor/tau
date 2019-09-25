# Tau

[![CircleCI](https://circleci.com/gh/avinor/tau.svg?style=svg)](https://circleci.com/gh/avinor/tau)
[![Go Report Card](https://goreportcard.com/badge/github.com/avinor/tau)](https://goreportcard.com/report/github.com/avinor/tau)
[![GoDoc](https://godoc.org/github.com/avinor/tau?status.svg)](https://godoc.org/github.com/avinor/tau)

Tau (Terraform Avinor Utility) is a thin wrapper for [Terraform](https://www.terraform.io/) that makes terraform execution and secret handling easier. It tries to conform to the principal of keeping code DRY (Don't Repeat Yourself) and bases all executions on terraform modules. Since terraform already provides a lot of excellent features it tries to not change that too much but use those features best possible.

## Highlights

**DRY:** All executions are based on using terraform modules. Tau configuration file just describes how to execute / deploy the modules.

**Dependency handling:** Tau handles dependencies between different modules. Executing them in correct order and making output from one module available to others.

**Secret handling:** Recommended way to deal with secrets is passing them in as input variables. Using the build-in data sources in terraform tau can send anything as input variables to terraform.

**Backend:** Backend configuration is taken out of modules and defined in tau configuration.

## Installation

1. Tau requires [terraform](https://www.terraform.io/) 0.12+, [download](https://www.terraform.io/downloads.html) and install first
2. Download tau from [Release page](https://github.com/avinor/tau/releases) for your OS
3. Rename file to `tau` and add it to your `PATH`

Alternatively clone this repository and run `go install` to install latest version.

## How it works

1. Create a new module in terraform, or use an existing one. Lets use `avinor/kubernetes/azurerm` as an example.

2. Create a tau configuration file (ending in `.hcl` or `.tau`) that includes the module and sets input parameters

```terraform
module {
    source = "avinor/kubernetes/azurerm"
    version = "1.0.0"
}

inputs {
    name = "example"
    resource_group_name = "example-rg"
    service_cidr = "10.241.0.0/24"
    kubernetes_version = "1.13.5"

    ... node pools etc ...
}
```

3. Running `tau init` will create a new directory called `.tau` where it downloads the module and then runs `terraform init`.

4. Run `tau plan` and it will create a `terraform.tfvars` file in module directory with the parsed input values before running `terraform plan`.

So far it has just wrapped terraform to do some of the configurations a different way. Lets use some of the more powerful features of tau.

5. Most likely we want state to be stored in a remote storage. Define a backend block in same configuration file to define backend.

```terraform
backend "azurerm" {
    storage_account_name = "tfstate"
    container_name       = "state"
    key                  = "kubernetes.tfstate"
}
```

6. Running `tau init` will now in additon to downloading module also create a `tau_override.tf` file in module directory that configures the backend.

7. If using multiple subscriptions and environments we might want to ensure this module is deployed in correct environment. Define an `environment_variables` block to set some env vars for terraform commands.

```terraform
environment_variables {
    ARM_SUBSCRIPTION_ID = "xxxx-xxxx-xxxx-xxxx"
}
```

8. Some of the input variables might depend on output from another module / tau deployment. By using a dependency it can pass those outputs and make them available for input variables.

```terraform
dependency "vnet" {
    source = "./vnet.hcl"
}

inputs {
    ...

    agent_pools = [
        {
            name = "ipt"
            vm_size = "Standard_D2_v3"
            vnet_subnet_id = dependency.vnet.outputs.subnets.aks
        },
    ]
}
```

9. When running `tau plan` now it will try to resolve the dependencies first by reading the output variables from dependencies remote state. By using the remote state it only needs access to the state file for dependency, and not require to execute any terraform commands. Once read from remote state it will use that as input variable when running `terraform plan`.

10. In addition to reading from another module output it might be necessary to read secrets too. Any `data` blocks defined in configuration file will be resolved first and sent as input to module.

```terraform
data "azurerm_key_vault_secret" "sp" {
    name = "sp-secret"
    key_vault_id = "/subscriptions/xxxx-xxxx-xxxx-xxxx/resourceGroups/secrets-rg/providers/Microsoft.KeyVault/vaults/terraform-secrets-kv"
}

inputs {
    ...

    service_principal = {
        client_id = data.azurerm_key_vault_secret.secret.value
        client_secret = data.azurerm_key_vault_secret.secret.value
    }
}
```

11. Running `tau plan` will in addition to creating temporary dependency module now also create a temporary data module defining all `data` blocks from configuration. Output from those data blocks can be used in input variables like in normal terraform code.

12. All this have created an almost complete configuration. As a last step we want it to perform some initialization first to setup account access. Create a prepare hook to run initialization.

```terraform
hook "set_access_key" {
    trigger_on = "prepare"
    command = "./set_access_key.sh"
    set_env = true
}
```

Usually we are not deploying only one module but many. When running tau commands it will by default process all files in same folder, except those where filename ends in `_auto`. These are merged together in all deployment files.

13. To share some of the initialization between all modules in same folder create a new file `common_auto.hcl` in same folder and move the `hook` and `environment_variables` blocks to new file.

```terraform
environment_variables {
    ARM_SUBSCRIPTION_ID = "xxxx-xxxx-xxxx-xxxx"
}

hook "set_access_key" {
    trigger_on = "prepare"
    command = "./set_access_key.sh"
    set_env = true
}
```

It is now configured to execute hook and set environment variables for all executions in same folder. Any new configuration files in same folder will execute same hook.

## Remote state

See documentation on how to handle [single](./docs/single_backend.md) and [multiple backends](./docs/multiple_backends.md).

## Configuration

Any files named `.hcl` or `.tau` are read, where each file is one deployment of module. Based on the example in "How it works" section above it could end up like this:

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

    // override backend config from logs.hcl
    backend {
        sas_token = "override"
    }
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

```terraform
hook "set_access_key" {
    #Event to trigger hook on. Possible values are "prepare" and "finish"
    trigger_on = "prepare"

    # Command to execute, should not include arguments
    command = "az"

    # Alternative to defining command, reference to script to execute
    script = "https://raw.githubusercontent.com/avinor/tau/master/hack/az_copy_output_from_state.sh"

    # Arguments to send to command
    args = ["aks", "get-credentials"]

    # If true it will read output in format "key = value"
    set_env = true

    # Fail on error or continue running ignoring error
    fail_on_error = false

    # Disable cache and make sure command is run every time
    disable_cache = false

    # Working directory when executing command
    working_dir = "/tmp"
}
```

One or more hooks that triggers on specific events during deployment. It can read the output from command run and set environment variables for terraform, for instance access keys etc. `trigger_on` defines which event to trigger the hook on. This can either be just simple event (`prepare` or `finish`) or it can include which commands to trigger for. If hook should only trigger on `init` command, but not any other, then define `trigger_on` as `prepare:init`. Arguments after : is a comma separate list of commands to execute on.

Either `command` or `script` has to be defined. A command can be any locally available command, or local script, while a script is retrieved by using go-getter and can therefore be a script in a remote git repository as well. See [go-getter](https://github.com/hashicorp/go-getter) for download options.

To read output and set environment variables set `set_env` = true. It will read all output in format "key = value" and add them to the environment when running terraform.

If `fail_on_error` is set it will accept any failures from command and continue executing terraform commands. Default value is false and it will stop all executions.

To optimize execution and not run same command multiple times (for instance retrieving same access key) it caches output from every command and reuses cached value if called multiple times in same run. To disable cache set `disable_cache` = true.

### dependency

```terraform
dependency "logs" {
    # Source of dependency, has to be a local file
    source = "./logs.hcl"

    # Resolve the dependency in separate environment
    run_in_separate_env = true

    # Override one or all of attributes from dependency backend configuration
    backend {
        sas_token = "override"
    }
}
```

One or more dependencies for this deployment. Using dependency block has 2 effects:

* Make sure modules are deployed in correct order
* Make output from dependency available as variables

When resolving the output from a dependency it does this by using the terraform remote_state data source. Using example above it has a dependency on vnet.hcl that provides an output map of all subnets with their ids. Tau will not try to run any of the dependencies as that could require access it does not have, for instance vnet could be deployed in another subscription. Instead it creates a temporary terraform script that defines one `terraform_remote_state` data source for each variable defined in input block. It reads the backend definition from dependency source, but backend configuration can be overriden with the backend block in dependency definition. By doing it this way it should not be necessary to define any `terraform_remote_state` inside the module itself, and reading output from another module only requires access to its state store.

By default it will inherit the same environment variables (from hooks as well) as current deployment, unless `run_in_separate_env` attribute is set to true. When this is set to true it will not inherit any environment variables and that dependency will be resolved by running any hooks defined in dependency first. This is useful if dependency is deployed in different subscription.

### data

Data can be any data source available in terraform. This could be used to read secrets from a key vault, get Kubernetes versions etc. These will be resolved in same context as module is running, with same environment variables.

Output from data source can be used same way as in terraform by using the `data.source...` variables.

See terraform documentation for configuration of data blocks.

### environment_variables

```terraform
environment_variables {
    # Any key = value pair of environment variables
    ARM_SUBSCRIPTION_ID = "xxxx-xxxx-xxxx-xxxx"
}
```

A list of key value pair of environment variables that should be added to the context before running any terraform commands. This could be access keys, subscription ids etc. It is also possible to set environment variables from [hooks](#hook).

### backend

```terraform
backend "azurerm" {
    # Configuration for the specific backend
    ...
}
```

Backend to use for remote state storage. The configuration is same as in terraform, so look in terraform documentation how to configure this for each available remote backend.

Tau will create an override file with backend definition before running the module. By doing this it is not required to define any backend configuration in the module.

### module

```terraform
module {
    # Source of dependency, supports any go-getter + terraform registry
    source = "avinor/kubernetes/azurerm"

    # Terraform registry version, if using terraform registry
    version = "1.0.0"
}
```

Module is the source, and optionally version, of module to deploy. Source can be any sources available in go-getter library (http(s), git, local file, s3...) and terraform registry. If the version attribute is defined it will assume that source is from a terraform registry and will attempt to download from registry.

### inputs

Variable inputs to send to module on execution. Can contain references to any data source and dependencies. Before executing plan / apply it will create a `terraform.tfvars` file in the module temporary folder with all resolved variables. It is important to remember that even secrets sent as input variables are stored in remote state.

## Variables

In addition to the `data.` and `dependency.` variables that are resolved by terraform there are some predefined variables available. In this context source is the configuration file that is currently being processed. When reading included files the source variable will be origin file, not file that is included.

Example column shows result based on file `/tmp/virtual-network.hcl`

variable | Description | Example
---------|-------------|--------
source.path     | Full path for source file | /tmp/virtual-network.hcl
source.name     | Name of source file without extension | virtual-network
source.filename | Filename of source file, same as name just with extension | virtual-network.hcl
module.path     | Path where module will be downloaded, might not exist early in execution | /tmp/virtual-network.hcl/.tau/virtual-network.hcl/module

Variables can be used when defining backend configuration in auto imported files for instance. By using `source.name` it will resolve to name of source file during processing.

## Auto import

When executing a file or folder it will by default ignore all files ending in `_auto.(hcl|tau)` as those are considered auto import files. It will instead merge those files together with source file. Auto files can be used to define common settings across all modules in same folder. Using variables in auto files makes it possible to define a common backend configuration that will change based on source file being executed.

`common_auto.hcl`

```terraform
backend "azurerm" {
    storage_account_name = "tfstate"
    container_name       = "state"
    key                  = "westeurope/${source.name}.tfstate"
}
```

`virtual-network.hcl`

```terraform
module {
    ...
}

inputs {
    ...
}
```

When running `tau init -f virtual-network-hcl` it will load the `common_auto.hcl` file first and replace `{source.name}` with `virtual-network` since that is the source file. Then it will merge configuration with that from `virtual-network.hcl` file.

## CI Pipeline

When using terraform in a CI pipeline it is recommended to first run plan, then have manual approval of some sort of the plan before running apply. To keep the same plan files from plan stage the entire `.tau` directory can be saved between the stages. Restoring the directory into same folder in apply stage it is possible to run `tau apply` directory to apply all changes from plan.

## Comparison

There are other great tools for deploying terraform modules as well. This is a short comparison of them and why we wrote tau.

### Terragrunt

[Terragrunt](https://github.com/gruntwork-io/terragrunt) is a nice tool to reuse modules for multiple deployments, however we disagreed on some of the design choices and also wanted to have a better way to handle secrets.

We started using terragrunt but with the 0.12 release of terraform it stopped working until terragrunt updated to latest terraform syntax. While we had this challenge we also had some issues on how to handle secrets and needed a better way for that.

#### Backend defined in module

Terragrunt requires that the backend is defined in the module because it uses `terraform init -from-module ...`, so backend block has to exist to override.

Tau instead uses `go-getter` directly to download the module and then creates an override file defining the backend. It will run `terraform init` afterwards when all overrides have been defined.

#### Dependencies

Terragrunt recommends using a `terraform_remote_state` data source to retrieve dependencies in modules. We disagree that this should be in the module and rather a part of input. See documentation for more details how this is solved in tau.

#### Secrets

Terragrunt don't have any way to handle secrets. It supports sending environment variables to commands, but cannot retrieve secrets from an Azure Key Vault for instance.

### Astro

[Astro]((https://github.com/uber/astro/)) is a tool created by Uber to manage terraform executions.

#### YAML

There is nothing wrong with yaml format, but terraform has decided to use its own format called hcl. Since tool is working on terraform code we find it more natural to also use hcl instead of yaml.

#### One big file

Astro defines all modules in one big file instead of supporting separate files for each deployment.

#### Secrets

No secret handling, does not support retrieving values from a secret store.
