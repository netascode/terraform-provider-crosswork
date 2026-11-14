resource "crosswork_golden_config_templates" "example" {
  name           = "cwtest-ntp"
  description    = "NTP servers"
  tags           = ["core"]
  vendor         = "Cisco Systems"
  software_type  = "IOS XR"
  product_series = ["Cisco XRv 9000 Series Virtual Routers"]
  config         = file("${path.module}/cwtest-ntp.j2")
  variables = [
    {
      name  = "ntp_server"
      value = "10.1.1.2"
    },
    {
      name             = "ntp_key"
      value_wo         = var.ntp_key
      value_wo_version = 1
    },
  ]
}
