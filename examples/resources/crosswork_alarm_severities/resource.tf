resource "crosswork_alarm_severities" "example" {
  severities = [
    {
      name               = "L2-ELMI-4-AIPC"
      severity           = "minor"
      auto_clear_minutes = 1440
    }
  ]
}
