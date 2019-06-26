[![Go Report Card](https://goreportcard.com/badge/github.com/avinor/tau)](https://goreportcard.com/report/github.com/avinor/tau)
[![GoDoc](https://godoc.org/github.com/avinor/tau?status.svg)](https://godoc.org/github.com/avinor/tau)

# Tau

Tau (Terraform Avinor Utility) is a thin wrapper for [Terraform](https://www.terraform.io/) that makes deployment and secret handling easier. The basic idea is taken from [Terragrunt](https://github.com/gruntwork-io/terragrunt) to keep code DRY.

## Installation

- Download and have fun!

## Compared to terragrunt

Terragrunt is a nice tool to reuse modules for multiple deployments, however we disagreed on some of the design choices and also wanted to have a better way to handle secrets.

- backend in module
- data sources
