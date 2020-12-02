dependency "logs" {
  source = "./logs.hcl"
}

module {
  source  = "avinor/virtual-network-hub/azurerm"
  version = "1.3.0"
}

inputs {
  name                       = "simple"
  resource_group_name        = "networking-simple-rg"
  location                   = "westeurope"
  address_space              = "10.0.0.0/24"
  log_analytics_workspace_id = dependency.logs.outputs.resource_id

  public_ip_names = [
    "simple"
  ]

  management_nsg_rules = [
    {
      name                       = "allow-ssh"
      direction                  = "Inbound"
      access                     = "Allow"
      protocol                   = "Tcp"
      source_port_range          = "*"
      destination_port_range     = "22"
      source_address_prefix      = "VirtualNetwork"
      destination_address_prefix = "VirtualNetwork"
    },
  ]

  firewall_application_rules = [
    {
      name             = "linux"
      action           = "Allow"
      source_addresses = ["10.0.0.0/8"]
      target_fqdns = [
        "*.ubuntu.com",
        "*.snapcraft.io",
        "*.opensuse.org",
      ]
      protocol = {
        type = "Https"
        port = "443"
      }
    },
  ]

  firewall_network_rules = [
    {
      name                  = "ntp"
      action                = "Allow"
      source_addresses      = ["10.0.0.0/8"]
      destination_ports     = ["123"]
      destination_addresses = ["*"]
      protocols             = ["UDP"]
    },
  ]
}