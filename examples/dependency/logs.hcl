module {
  source  = "avinor/log-analytics/azurerm"
  version = "1.1.0"
}

inputs {
  name                = "shared-logs"
  resource_group_name = "logs-rg"
  location            = "westeurope"
  sku                 = "PerNode"

  solutions = [
    {
      solution_name = "ContainerInsights",
      publisher     = "Microsoft",
      product       = "OMSGallery/ContainerInsights",
    },
  ]
}