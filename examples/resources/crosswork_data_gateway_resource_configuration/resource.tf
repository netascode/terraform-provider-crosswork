resource "crosswork_data_gateway_resource_configuration" "example" {
  collectors = [
    {
      name    = "CLI"
      memory  = 9217
      enabled = true
    }
  ]
}
