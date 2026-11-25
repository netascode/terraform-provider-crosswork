resource "crosswork_device_alarm_manager" "example" {
  models = { "Cisco ASR 9010 Router" = true, "Cisco ASR 9006 Router" = false }
}
