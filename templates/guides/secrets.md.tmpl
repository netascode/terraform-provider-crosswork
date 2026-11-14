---
subcategory: "Guides"
page_title: "Managing Secrets"
description: |-
  Configure and rotate write-only secrets without storing their values in Terraform plans or state.
---

# Managing Secrets

Use Terraform 1.11 or later for [write-only arguments](https://developer.hashicorp.com/terraform/language/manage-sensitive-data/write-only). The provider exposes these arguments with a `_wo` suffix and a corresponding `_wo_version` argument. Terraform keeps the version in state, but omits the write-only value from saved plans and state.

The older secret arguments, such as `password`, remain supported. They are marked sensitive to redact normal Terraform output, but their values are still stored in state. Use the write-only alternatives when the value should not be persisted.

## Credential profile example

Declare secret input variables as both `sensitive` and `ephemeral`. The ephemeral setting also prevents Terraform from storing the input variable values in saved plans; using only `sensitive` does not provide that protection.

```terraform
variable "ssh_password" {
  type      = string
  sensitive = true
  ephemeral = true
}

variable "enable_password" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "crosswork_credential_profiles" "example" {
  profile = "device-ssh"
  credentials = [
    {
      type                             = "ROBOT_USERPASS_SSH"
      user_name                        = "network-admin"
      password_wo                      = var.ssh_password
      password_wo_version              = 1
      enable_password                  = "ROBOT_USERPASSEN_ENABLE"
      enable_password_value_wo         = var.enable_password
      enable_password_value_wo_version = 1
    }
  ]
}
```

Supply these inputs at runtime, for example through the `TF_VAR_ssh_password` and `TF_VAR_enable_password` environment variables. Supply them again when applying a saved plan, because Terraform does not retain ephemeral inputs in that plan. Write-only arguments do not protect secrets stored in configuration files or other secret sources.

## Rotating secrets

Each write-only secret has its own version. When you change `ssh_password`, increment `password_wo_version`, for example from `1` to `2`. When you change the enable password, increment `enable_password_value_wo_version`. The version is an integer that asks Terraform for an update; it is not a secret and it is not sent to Crosswork.

Changing only the write-only value produces no Terraform diff, because Terraform has no previous value to compare. Increment the version to send the new value.

Keep every write-only value available whenever you apply an update to the resource, including an update to an ordinary attribute such as `user_name`. The provider sends the whole record on each update and reads the write-only values from the current configuration. A credential profile updated without its password is refused by Crosswork with `Password cannot be empty`.

The provider cannot recover write-only values from state, and it cannot detect a secret that was changed outside Terraform.

## Resources with write-only secrets

The same argument and version pattern is available for these resources:

| Resource | Write-only secret arguments |
| --- | --- |
| `crosswork_users` | `password_wo` |
| `crosswork_credential_profiles` | `password_wo`, `enable_password_value_wo` (both inside `credentials`), `snmpv2_read_community_wo`, `snmpv2_write_community_wo`, `snmpv3_auth_password_wo`, `snmpv3_priv_password_wo` |
| `crosswork_backup_destination` | `password_wo` |
| `crosswork_trap_destination` | `community_wo`, `auth_password_wo`, `privacy_password_wo` |
| `crosswork_golden_config_global_variables` | `value_wo` |
| `crosswork_golden_config_templates` | `value_wo` of a variable in `variables` |
| `crosswork_golden_config_jobs` | `value_wo` of a variable of a device in `devices` |

For each argument in the table, its version argument adds `_version`, for example `snmpv3_auth_password_wo_version`. A version is required when its write-only argument is configured. Use either the older secret argument or its write-only alternative, never both for the same secret. Consult each resource's schema for its exact arguments.

## Values that Crosswork returns

Crosswork does not treat every secret the same way on a read:

- **Masked secrets:** the backup destination returns its password masked as `********` and needs the real password on every save. The provider keeps the value of the configuration and never reads a secret back from Crosswork.
- **Variables in clear text:** Golden Configuration returns the value of a variable in clear text on every read, and the status of a job repeats its input. Use `value_wo` for a secret variable: it keeps the value out of the plan and the state of Terraform, but anyone who may read the variable in Crosswork can still see it.

## Golden Configuration variables

A Golden Configuration variable has `value` for text and `value_wo` for a secret, and exactly one of them must be set. The secret is raised with `value_wo_version`:

```terraform
variable "ntp_key" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "crosswork_golden_config_global_variables" "secret" {
  name             = "ntp_key"
  value_wo         = var.ntp_key
  value_wo_version = 1
}
```

A template variable uses the same pair inside the `variables` list of the template, and a job uses it inside the `variables` of a device in `devices`.

## Migrating existing secrets

An upgrade alone does not require switching existing secret arguments. To migrate a secret, keep the address of the resource, remove the older argument such as `password`, and configure `password_wo` with the current secret plus `password_wo_version = 1`. Review and apply the plan; Terraform tracks the version and removes the old secret from the current state. Other resource-specific changes can still cause a replacement, as the schema of the resource describes.

Migration does not remove values from older state snapshots, backups or saved plans. Those files still contain any secrets that were stored in ordinary or sensitive arguments.
