[![CircleCI](https://circleci.com/gh/avinor/tau.svg?style=svg)](https://circleci.com/gh/avinor/tau)
[![Go Report Card](https://goreportcard.com/badge/github.com/avinor/tau)](https://goreportcard.com/report/github.com/avinor/tau)
[![GoDoc](https://godoc.org/github.com/avinor/tau?status.svg)](https://godoc.org/github.com/avinor/tau)

# Tau

Tau (Terraform Avinor Utility) is a thin wrapper for [Terraform](https://www.terraform.io/) that makes terraform execution and secret handling easier. The basic idea is taken from [Terragrunt](https://github.com/gruntwork-io/terragrunt) to keep code DRY, however it does some changes to how remote state and data sources are handled. It tries to not change terraform too much and use the existing functionallities. Hopefully that will make it easier to upgrade when terraform does major changes as well.

For a comparison with other similar tools see [comparison](#comparison).

**NOTE: This is still in development.**

## Installation

1. Tau requires [terraform](https://www.terraform.io/) 0.12+, download and install first
2. Download tau from [Release page](https://github.com/avinor/tau/releases) for your OS
3. Rename file to `tau` and add it to your `PATH`

## Motivation

At the time of creation the tools available did not feel complete enough. We evaluated a few options (terragrunt, astro, "pure" terraform), but decided to create a new tool. See [comparison](#comparison) for details on why other tools where dismissed in the end.

Goals of project was to handle all these issues:

- Keep code DRY
- Backend configuration
- Secrets
- Dependencies

## Configuration

For a simple example how to get started create a file `example.hcl` (or `example.tau`) and include a module and some inputs:

```terraform
module {
    source = "avinor/storage-account/azurerm"
    version = "1.0.0"
}

inputs {
    name = "example"
    resource_group_name = "example-rg"
    location = "westeurope"

    containers = [
        {
            name = "example"
            access_type = "private"
        },
    ]
}
```

Run `tau init` in same folder and it will download module `avinor/storage-account/azurerm`, create input file and run `terraform init`.

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
