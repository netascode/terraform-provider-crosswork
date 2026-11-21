resource "crosswork_aaa_password_policy" "example" {
  min_password_length                  = 9
  no_username                          = true
  no_cisco_variant                     = true
  no_char_repetition                   = true
  num_changed_chars_enabled            = true
  num_changed_chars                    = 2
  num_reuse_limit_enabled              = true
  num_reuse_limit                      = 6
  password_reuse_days_enabled          = true
  password_reuse_days                  = 16
  failed_logins_before_lockout_enabled = true
  failed_logins_before_lockout         = 6
  lock_out_user_time_enabled           = true
  lock_out_user_time                   = 7
  password_expiry_days_enabled         = true
  password_expiry_days                 = 61
  days_for_warning_enabled             = true
  days_for_warning                     = 16
  change_password_on_first_login       = true
}
