resource "crosswork_golden_config_jobs" "example" {
  template_name      = crosswork_golden_config_templates.example.name
  template_config_id = crosswork_golden_config_templates.example.config_id
  vendor             = "Cisco Systems"
  product_series     = "Cisco XRv 9000 Series Virtual Routers"
  mop                = "Default Golden Config - Cisco"
  devices = {
    "build-router-1" = {}
    "build-router-2" = {
      variables = [{
        name  = "ntp_server"
        value = "10.2.2.2"
      }]
    }
  }
}
