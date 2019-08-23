backend "azurerm" {
    storage_account_name = "abcd1234"
    container_name       = "tfstate"
    key                  = "prod.terraform.tfstate"
    sas_token            = env("SAS-TOKEN")
}

module {
    source = "avinor/virtual-network-hub/azurerm"
    version = "1.0.0"
}

inputs {
    name = "simple"
    resource_group_name = "simple-rg"
    location = "westeurope"
    address_space = "10.0.0.0/22"
}