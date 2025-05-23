resource "forgejo_repository_actions_secret" {
  data       = "secret"
  name       = "TEST"
  owner      = "adyxax"
  repository = "example"
}
