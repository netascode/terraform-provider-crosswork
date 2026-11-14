resource "crosswork_golden_config_global_variables" "example" {
  name  = "cwtest_ntp_server"
  type  = "STRING"
  value = "10.1.1.1"
}

# A secret stays out of the plan and the state.
resource "crosswork_golden_config_global_variables" "secret" {
  name             = "cwtest_ntp_key"
  value_wo         = var.ntp_key
  value_wo_version = 1
}
