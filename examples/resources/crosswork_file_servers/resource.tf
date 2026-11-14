resource "crosswork_file_servers" "example" {
  ftp_enabled         = true
  sftp_upload_enabled = true
  tftp_enabled        = false
}
