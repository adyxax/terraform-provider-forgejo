resource "forgejo_repository_actions_variable" {
  data       = "value"
  name       = "test"
  owner      = "adyxax"
  repository = "example"
}
