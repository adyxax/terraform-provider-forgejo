data "forgejo_organizations" "main" {}

data "forgejo_teams" "main" {
  for_each = toset([for org in data.forgejo_organizations.main.elements :
    org.name
  ])

  organization_name = each.key
}
