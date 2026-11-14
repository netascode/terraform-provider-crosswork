resource "crosswork_provider_profiles" "example" {
  name    = "cwtest-provider-profiles"
  family  = "ROBOT_PROVIDER_NSO"
  profile = "cwtest-provider-profile"
  connectivity = [
    {
      protocol   = "ROBOT_MSVC_TRANS_SSH"
      ip_address = "192.0.2.50"
      port       = 22
      timeout    = 60
    }
  ]
  supported_models = [
    {
      model_prefix  = "Cisco-IOS-XR"
      model_version = "7.8.2"
    }
  ]
  properties = { nso_crosslaunch_url = "https://nso.example.invalid/nso" }
}
