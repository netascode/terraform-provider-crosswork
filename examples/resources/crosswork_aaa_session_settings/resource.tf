resource "crosswork_aaa_session_settings" "example" {
  fallback_type              = "FB_ON_NO_RESPONSE"
  idle_session_timeout       = 60
  parallel_sessions          = 100
  parallel_sessions_per_user = 100
  source_ip_auditing         = true
}
