resource "crosswork_l3vpn_source_prefix_sets" "example" {
  name = "cwtest-source-prefixes"
  mode = "ipv4"
  prefix_list = [
    {
      ip_prefix = "10.250.1.0/24"
    }
  ]
}
