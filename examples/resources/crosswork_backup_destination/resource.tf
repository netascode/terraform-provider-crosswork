resource "crosswork_backup_destination" "example" {
  host_name           = "192.0.2.10"
  port                = 22
  username            = "cwtest"
  password_wo         = "cwtest-password"
  password_wo_version = 1
  remote_path         = "/backups"
}
