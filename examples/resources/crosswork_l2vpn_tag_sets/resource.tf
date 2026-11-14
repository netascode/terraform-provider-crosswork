resource "crosswork_l2vpn_tag_sets" "example" {
  name      = "cwtest-l2-tags"
  tag_value = [4812]
  co_flag   = "10"
}
