resource "crosswork_aaa_sso_services" "example" {
  services = [
    {
      name              = "cwtestsso"
      evaluation_order  = 3
      metadata          = "PE1ldGFkYXRhLz4="
      metadata_location = "cwtest-metadata.xml"
    }
  ]
}
