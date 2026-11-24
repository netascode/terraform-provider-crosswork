resource "crosswork_device_notification_policies" "example" {
  name        = "cwtest-dnp"
  description = "Critical alarms"
  filters = [
    {
      severities = ["critical"]
    }
  ]
  trap_receivers = [
    {
      name = "cwtest-dnp-trap"
    }
  ]
  syslog_receivers = [
    {
      name = "cwtest-dnp-syslog"
    }
  ]
}
