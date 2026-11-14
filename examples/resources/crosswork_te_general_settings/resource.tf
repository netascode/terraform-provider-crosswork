resource "crosswork_te_general_settings" "example" {
  policy_unused_threshold                  = 2000
  policy_provisioning_timeout              = 120
  bandwidth_on_demand_provisioning_timeout = 100
  igp_path_request_timeout                 = 200
  policy_details_request_timeout           = 220
}
