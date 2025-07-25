resource "forgejo_repository" "user_example" {
  description = "test repository"
  name        = "test"
  private     = false
}

resource "forgejo_repository" "organization_example" {
  description = "test repository"
  name        = "test"
  owner       = "adyxax.org"
  private     = false
}
