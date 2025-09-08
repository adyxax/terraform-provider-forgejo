# Importing a repository belonging to the user whose credentials the provider
# was instantiated with:
terraform import forgejo_repository.main <repository_name>

# Importing a repository that belongs to an organization:
terraform import forgejo_repository.main <owner>/<repository_name>
