[![Go Report Card](https://goreportcard.com/badge/github.com/avinor/tau)](https://goreportcard.com/report/github.com/avinor/tau)
[![GoDoc](https://godoc.org/github.com/avinor/tau?status.svg)](https://godoc.org/github.com/avinor/tau)

# Tau

Tau (Terraform Avinor Utility) is a thin wrapper for [Terraform](https://www.terraform.io/) that makes deployment and secret handling easier. The basic idea is taken from [Terragrunt](https://github.com/gruntwork-io/terragrunt) to keep code DRY, however it does some changes to how remote state and data sources are handled. It tries to not change terraform to much and use the existing functionallities.

NOTE: This is still in development.

## Installation

1. Tau requires [terraform](https://www.terraform.io/) 0.12+, download and install first
2. Download tau from [Release page](https://github.com/avinor/tau/releases) for your OS
3. Rename file to `tau` and add it to your `PATH`

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

## Compared to terragrunt

Terragrunt is a nice tool to reuse modules for multiple deployments, however we disagreed on some of the design choices and also wanted to have a better way to handle secrets. 

- backend in module
- data sources
