resource "crosswork_srte_policies" "example" {
  name                                       = "cwtest-policy"
  head_end                                   = "build-router-1"
  tail_end                                   = "192.0.2.61"
  color                                      = 999
  binding_sid                                = 24001
  bandwidth                                  = 1000
  source_address                             = "192.0.2.1"
  delay_measurement_profile                  = "delay-profile-1"
  delay_measurement_logging_delay_exceeded   = true
  performance_measurement_reverse_path_label = 2000
  auto_route_metric_relative                 = -5
  auto_route_include_prefixes = [
    {
      prefix_address = "10.250.0.0/16"
    }
  ]
  auto_route_force_sr_include = true
  auto_route_forward_class    = 2
  paths = [
    {
      preference = 100
      explicit_sid_lists = [
        {
          name   = "cwtest-sid-list"
          weight = 10
        }
      ]
      explicit_affinity_rules = [
        {
          action = "include-any"
          colors = ["red"]
        }
      ]
      explicit_segments_sid_algorithm = 128
      explicit_disjoint_path_type     = "link"
      explicit_disjoint_path_group_id = 1
      explicit_disjoint_path_sub_id   = 1
    }
  ]
}
