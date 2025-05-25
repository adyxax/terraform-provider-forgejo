resource "forgejo_repository_actions_secret" "main" {
  data       = "secret"
  name       = "TEST"
  owner      = "adyxax"
  repository = "example"
}
