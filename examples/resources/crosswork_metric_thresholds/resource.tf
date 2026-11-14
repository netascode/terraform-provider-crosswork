resource "crosswork_metric_thresholds" "example" {
  thresholds = {
    "Utilization" = {
      levels = [
        {
          severity    = "CRITICAL"
          upper_value = 100
        }
      ]
    }
  }
}
