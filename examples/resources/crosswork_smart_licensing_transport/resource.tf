resource "crosswork_smart_licensing_transport" "example" {
  mode       = "DIRECT"
  url        = "https://smartreceiver.cisco.com/licservice/license"
  proxy_host = "192.0.2.1"
  proxy_port = 8080
}
