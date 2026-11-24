resource "crosswork_suppression_policies" "example" {
  name        = "cwtest-suppress"
  description = "Suppress the warm start alarms"
  action      = "suppressAlarm"
  criteria    = "eventType in [WARMSTART]"
}
