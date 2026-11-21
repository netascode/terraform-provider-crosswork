resource "crosswork_aaa_tacacs_servers" "example" {
  servers = [
    {
      host               = "192.0.2.10"
      port               = 49
      secret_wo          = "AAASecret"
      secret_wo_version  = 1
      secret_format      = "ascii"
      service            = "crosswork"
      priority           = 30
      retransmit_timeout = 10
      retries            = 3
      auth_type          = "pap"
      interface          = "192.0.2.1"
      policy_id          = "role0"
      virtual_domain     = "deviceAccessGroups"
    }
  ]
}
