resource "crosswork_collection_jobs" "example" {
  application_id = "cwtest-app"
  context_id     = "cwtest-context"
  collector_type = "CLI"
  devices        = ["build-router-1"]
  destination    = "cwtest-dest-kafka"
  cli_sensors = [
    {
      command = "show version"
      cadence = 60
      topic   = "cwtest-topic"
    }
  ]
}
