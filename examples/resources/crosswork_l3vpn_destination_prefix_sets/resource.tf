resource "crosswork_l3vpn_destination_prefix_sets" "example" {
  name = "cwtest-dest-prefixes"
  mode = "ipv4"
  prefix_list = [
    {
      ip_prefix = "10.250.1.0/24"
    }
  ]
}
