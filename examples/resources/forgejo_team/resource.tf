resource "forgejo_team" "main" {
  name              = "test"
  organization_name = "test"
  permission        = "read"
}
