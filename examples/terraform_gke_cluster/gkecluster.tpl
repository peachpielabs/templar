resource "google_service_account" "default" {
  account_id   = "{{.service_account_id}}"
  display_name = "{{.service_account_name}}"
}

resource "google_container_cluster" "primary" {
  name     = "{{.cluster_name}}"
  location = "{{.cluster_location}}"

  remove_default_node_pool = true
  initial_node_count       = {{.node_count}}
}
