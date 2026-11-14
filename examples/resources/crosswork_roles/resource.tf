resource "crosswork_roles" "example" {
  name = "cwtest-role"
  api_permissions = {
    "cw-fault-events-api" = "RW"
    "nso-fp-dep-mngr"     = "RWD"
  }
  task_permissions = ["id_export_audit_logs_access", "id_view_audit_logs_access", "id_nso_fp_deployment_management"]
}
