resource "crosswork_assign_ports_to_port_group" "example" {
  group_id = "16f9288d-ab96-4ea5-8593-94b30cbf00ee"
  ports = [
    {
      device    = "build-router-1"
      port_name = "GigabitEthernet0/0/0/0"
    }
  ]
}
