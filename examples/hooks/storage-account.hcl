hook "set_access_key" {
  trigger_on = "prepare"
  script     = "https://raw.githubusercontent.com/avinor/terraform-azurerm-remote-backend/master/set-access-keys.sh"
  args       = ["terraformname"]
  set_env    = true
}

module {
  source  = "avinor/storage-account/azurerm"
  version = "1.1.0"
}

inputs {
  name                = "simple"
  resource_group_name = "simple-rg"
  location            = "westeurope"

  containers = [
    {
      name        = "test"
      access_type = "private"
    },
  ]
}