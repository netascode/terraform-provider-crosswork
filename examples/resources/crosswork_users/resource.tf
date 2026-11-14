resource "crosswork_users" "example" {
  username            = "cwtest-gen-user"
  password_wo         = "CwtestPass123!"
  password_wo_version = 1
  policy_id           = "admin"
  first_name          = "Test"
  last_name           = "User"
  device_access_groups = [
    {
      uuid        = "b84e1a37-ef0d-4b79-984c-b162249d2e23"
      domain_name = "ALL-ACCESS"
    }
  ]
}
