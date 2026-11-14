resource "crosswork_srte_odn_templates" "example" {
  name  = "cwtest-odn"
  color = 4242
  head_end = [
    {
      name = "build-router-1"
    }
  ]
  maximum_sid_depth              = 8
  pce_group                      = "cwtest-pce-group"
  bandwidth                      = 1000
  source_address                 = "192.0.2.1"
  dynamic_pce                    = true
  dynamic_flex_alg               = 129
  dynamic_metric_type            = "igp"
  dynamic_metric_margin_absolute = 5
  dynamic_affinity_rules = [
    {
      action = "include-any"
      color  = ["red"]
    }
  ]
  dynamic_disjoint_path_type                 = "node"
  dynamic_disjoint_path_group_id             = 5
  dynamic_disjoint_path_sub_id               = 1
  delay_measurement_profile                  = "delay-profile-1"
  delay_measurement_logging_delay_exceeded   = true
  performance_measurement_reverse_path_label = 2000
}
