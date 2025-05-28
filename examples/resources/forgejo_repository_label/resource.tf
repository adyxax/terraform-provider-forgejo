resource "forgejo_repository_label" "main" {
  color       = "0000ff"
  description = "blue label"
  name        = "test"
  owner       = "adyxax"
  repository  = "example"
}
