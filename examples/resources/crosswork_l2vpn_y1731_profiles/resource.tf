resource "crosswork_l2vpn_y1731_profiles" "example" {
  name                       = "cwtest-y1731"
  type                       = "delay"
  schedule_interval          = 5
  schedule_duration          = "5"
  probe_type                 = "burst"
  probe_message_count        = 10
  probe_message_period       = 100
  probe_measurement_interval = "60"
  probe_frame_size           = 1000
  probe_priority             = 5
  delay_statistics = [
    {
      type = "delay-two-way"
    }
  ]
  bucket_size    = 1
  bucket_archive = 20
}
