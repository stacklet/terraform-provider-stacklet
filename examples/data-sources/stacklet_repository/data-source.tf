# Fetch a repository by URL
data "stacklet_repository" "by_url" {
  url = "https://github.com/example/repo"
}

# Fetch a repository by UUID
data "stacklet_repository" "by_uuid" {
  uuid = "12345678-90ab-cdef-1234-567890abcdef"
}
