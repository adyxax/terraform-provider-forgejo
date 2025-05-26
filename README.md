# Forgejo terraform provider

[Read the full documentation
here](https://registry.terraform.io/providers/adyxax/forgejo/latest/docs)

The [Forgejo](https://forgejo.org/) terraform / OpenTofu provider is used to
interact with the resources supported by Forgejo. The provider needs to be
configured with the proper credentials before it can be used. It requires
terraform 1.0 or later.

## Example Usage

```hcl
terraform {
  required_providers {
    forgejo = {
      source = "adyxax/forgejo"
    }
  }
}

provider "forgejo" {
  api_token  = var.forgejo_api_token
  base_url = "https://git.adyxax.org/"
}
```
