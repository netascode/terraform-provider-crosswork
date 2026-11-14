resource "crosswork_health_settings" "example" {
  rows = {
    "CPU_cpuUtilization" = {
      unit = "PERCENTAGE"
      categories = [
        {
          level  = "HEALTHY"
          min    = 0
          max    = 25
          values = [[2]]
        }
      ]
    }
  }
}
