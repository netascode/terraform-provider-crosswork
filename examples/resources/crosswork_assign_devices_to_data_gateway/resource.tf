resource "crosswork_assign_devices_to_data_gateway" "example" {
  data_gateway_id = "8551f8dc-91eb-43ee-ad4f-6e1993ca524c"
  devices         = ["build-router-1"]
}
