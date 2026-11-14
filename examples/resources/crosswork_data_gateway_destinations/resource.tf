resource "crosswork_data_gateway_destinations" "example" {
  name               = "cwtest-gen-destination"
  data_source        = "data-gateway"
  server_type        = "kafka"
  encoding           = "gpbkv"
  compression_method = "snappy"
  max_message_size   = 100000000
  buffer_memory      = 52428800
  batch_size         = 1048576
  linger_time        = 2000
  request_timeout    = 30
  collectors = [
    {
      name            = "CLI"
      buffer_memory   = 104857600
      batch_size      = 500000
      linger_time     = 1000
      request_timeout = 45
    }
  ]
  connections = [
    {
      address       = "10.14.0.150"
      prefix_length = 20
      port          = 9192
    }
  ]
  secure = false
}
