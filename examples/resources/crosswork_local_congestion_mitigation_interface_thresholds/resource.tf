resource "crosswork_local_congestion_mitigation_interface_thresholds" "example" {
  domain_id = "32"
  interface_thresholds = [
    {
      node                  = "build-router-1"
      interface             = "GigabitEthernet0/0/0/14"
      utilization_threshold = 75
    }
  ]
}
