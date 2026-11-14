resource "crosswork_port_groups_rules" "example" {
  group_id = "16f9288d-ab96-4ea5-8593-94b30cbf00ee"
  name     = "cwtest-port-group"
  conditions = [
    {
      attribute = "name"
      operator  = "starts_with"
      value     = "Four"
    }
  ]
}
