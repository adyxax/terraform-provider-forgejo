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
