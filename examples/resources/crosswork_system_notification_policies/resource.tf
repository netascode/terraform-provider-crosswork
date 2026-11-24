resource "crosswork_system_notification_policies" "example" {
  name        = "cwtest-policy"
  description = "Forward all alarms"
  criteria    = "*"
  trap_receivers = [
    {
      name = "cwtest-anp-trap"
    }
  ]
  syslog_receivers = [
    {
      name = "cwtest-anp-syslog"
    }
  ]
}
