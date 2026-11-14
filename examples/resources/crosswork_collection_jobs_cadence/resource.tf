resource "crosswork_collection_jobs_cadence" "example" {
  reachability_check = 600
  snmp_check         = 1200
  te_tunnel_id       = 1200
  show_clock         = 1800
}
