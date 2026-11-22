resource "crosswork_performance_policies" "example" {
  name            = "cwtest-policy"
  policy_template = "deviceHealth"
  description     = "Test policy"
  contact         = "noc@example.com"
  active          = true
  device_groups   = ["2661ed2d-11fb-49c4-9f9f-f71f1c12767b"]
  schema_intervals = [
    {
      schema  = "CPU"
      seconds = 600
    }
  ]
  thresholds = {
    "cpuUtilization" = {
      conditions = [
        {
          relation = "GE"
          value    = "90"
          count    = 2
          severity = "CRITICAL"
        }
      ]
    }
  }
}
