resource "crosswork_aaa_radius_servers" "example" {
  servers = [
    {
      host               = "192.0.2.20"
      port               = 1812
      secret_wo          = "AAASecret"
      secret_wo_version  = 1
      secret_format      = "ascii"
      retransmit_timeout = 5
      retries            = 3
      auth_type          = "pap"
      priority           = 40
      policy_id          = "role0"
      virtual_domain     = "deviceAccessGroups"
    }
  ]
}
