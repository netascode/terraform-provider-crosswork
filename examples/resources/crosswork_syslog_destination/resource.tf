resource "crosswork_syslog_destination" "example" {
  name     = "cwtest-syslog"
  address  = "192.0.2.1"
  port     = 1514
  protocol = "UDP"
}
