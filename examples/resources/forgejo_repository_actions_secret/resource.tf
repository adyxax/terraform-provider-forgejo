resource "forgejo_repository_actions_secret" {
  data       = "secret"
  name       = "test"
  owner      = "adyxax"
  repository = "example"
}
