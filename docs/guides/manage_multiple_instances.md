---
subcategory: "Guides"
page_title: "Manage Multiple Instances"
description: |-
    How to manage multiple Crosswork clusters.
---

# Manage Multiple Instances

When it comes to managing multiple Crosswork clusters, one can create multiple provider configurations and distinguish them by `alias` ([documentation](https://www.terraform.io/language/providers/configuration#alias-multiple-provider-configurations)).

```terraform
provider "crosswork" {
  alias    = "cluster-dev"
  username = "user1"
  password = "Cisco123"
  url      = "https://10.1.1.1:30603"
}

provider "crosswork" {
  alias    = "cluster-prod"
  username = "user1"
  password = "Cisco123"
  url      = "https://10.1.1.2:30603"
}
```

The disadvantage here is that the `provider` attribute of resources cannot be dynamic and therefore cannot be used in combination with `for_each`, as an example. The issue is being tracked [here](https://github.com/hashicorp/terraform/issues/24476).

This provider offers an alternative approach where multiple clusters can be managed by a single provider configuration. The optional `instance` attribute, which is available in every resource and data source, selects the respective cluster by its name. This assumes that every cluster uses the same credentials.

```terraform
locals {
  instances = [
    {
      name = "cluster-dev"
      url  = "https://10.1.1.1:30603"
    },
    {
      name = "cluster-prod"
      url  = "https://10.1.1.2:30603"
    },
  ]
}

provider "crosswork" {
  username  = "admin"
  password  = "Cisco123"
  instances = local.instances
}

# Every resource and data source accepts the instance attribute.
resource "crosswork_tags" "example" {
  for_each = toset([for inst in local.instances : inst.name])
  instance = each.key
  name     = "example"
  category = "example"
}
```

## Notes

- **Default instance.** A resource or data source without `instance` uses the default instance: the cluster of the provider `url` (or `CW_URL`). If `url` is not set, the first entry of `instances` is the default.
- **Unknown names.** An `instance` value that is not listed in `instances` fails with the error "Instance does not exist in provider configuration".
- **Sessions.** The provider logs in once per cluster and process, however many resources use it, and logs out when the process ends. Two instances with the same URL share one session.
- **Imports.** An import identifier cannot carry the instance, so `terraform import` addresses the default instance. To import an object of another cluster, make that cluster the default of a provider configuration (use `alias` or reorder `instances` when `url` is not set).
- **Changing `instance`.** Changing `instance` on an existing resource does not force a replacement: the next apply sends the update to the newly selected cluster. Destroy and recreate the resource to move it between clusters.
