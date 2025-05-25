resource "forgejo_repository_actions_variable" "main" {
  data       = "value"
  name       = "test"
  owner      = "adyxax"
  repository = "example"
}
