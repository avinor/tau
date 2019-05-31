data "azurerm_key_vault_secret" "test" {
  name      = "secret-sauce"
  key_vault_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/vault1"
}

dependency "vnet" {
    source = "./virtual-network"
}

backend "azurerm" {
    storage_account_name = "abcd1234"
    container_name       = "tfstate"
    key                  = "prod.terraform.tfstate"
    sas_token            = env("SAS-TOKEN")
}

module {
    source = "avinor/storage-account/azurerm"
    version = "1.0.0"
}

inputs {
    name = "simple"
    resource_group_name = "simple-rg"
    location = "westeurope"

    containers = [
        {
            name = azurerm_key_vault_secret.test.value
            access_type = "private"
        },
        {
            name = dependency.vnet.name
            access_type = "private"
        },
    ]
}