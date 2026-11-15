resource "crosswork_provider_nso_policies" "example" {
  provider_name         = "nso"
  match                 = true
  match_rule            = "*"
  onboard_to_nso        = true
  onboard_to_nso_rule   = "*"
  onboard_from_nso      = false
  onboard_from_nso_rule = "*"
  sync_from_nso         = true
  sync_from_nso_rule    = "*"
  check_sync            = true
  check_sync_rule       = "*"
  neds = [
    {
      ned  = "cisco-iosxr-cli-7.76"
      rule = "product_info.software_type = 'IOS XR*'"
    }
  ]
}
