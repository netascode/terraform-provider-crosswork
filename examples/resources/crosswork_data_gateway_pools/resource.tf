resource "crosswork_data_gateway_pools" "example" {
  name                  = "cwtest-gen-pool"
  description           = "Pool of the production gateways"
  virtual_ips           = ["172.16.34.198"]
  prefix_length         = 24
  spares                = 0
  instances             = ["2d009d08-e1eb-46b8-a137-1ef15e20ec72"]
  cdg_readiness_timeout = 60
}
