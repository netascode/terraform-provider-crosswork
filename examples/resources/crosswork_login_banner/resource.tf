resource "crosswork_login_banner" "example" {
  show_message = true
  user_ack     = false
  message      = "Authorized use only"
  icon         = "BANNER_ICON_IMPORTANT"
  title        = "Notice"
}
