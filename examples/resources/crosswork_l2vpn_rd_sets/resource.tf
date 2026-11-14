resource "crosswork_l2vpn_rd_sets" "example" {
  name = "cwtest-l2-rds"
  rd   = ["65000:1"]
}
