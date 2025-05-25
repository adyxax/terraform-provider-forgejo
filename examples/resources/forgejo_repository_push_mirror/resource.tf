resource "forgejo_repository_push_mirror" "main" {
  owner           = "adyxax"
  remote_address  = "https://github.com/adyxax/tfstated"
  remote_password = "secret"
  remote_username = "adyxax"
  repository      = "example"
}
