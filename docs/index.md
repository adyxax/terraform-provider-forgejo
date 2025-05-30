---
layout: ""
page_title: "Provider: forgejo"
description: "The Forgejo provider allows you to interact with the Forgejo API."
---

# Forgejo Provider

The [Forgejo](https://forgejo.org/) provider is used to interact with the
resources supported by Forgejo. The provider needs to be configured with the
proper credentials before it can be used. It requires terraform 1.0 or later.

Use the navigation to the left to read about the available resources.

## Example Usage

```terraform
terraform {
  required_providers {
    forgejo = {
      source = "adyxax/forgejo"
    }
  }
}

provider "forgejo" {
  api_token = var.forgejo_api_token
  endpoint  = "https://git.adyxax.org/"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `base_uri` (String) Forgejo's HTTP base URI.

### Optional

- `api_token` (String, Sensitive) Forgejo's api token. If not defined, the content of the environment variable `FORGEJO_API_TOKEN` will be used instead.
