resource "google_compute_firewall" "{{.rule_name}}" {
  name    = "{{.rule_name}}"
  network = google_compute_network.default.name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  {{$first_tag := true}}
  source_tags = [
    {{ range $tag := .source_tags }}{{if $first_tag}}{{$first_tag = false}}{{else}},
    {{end}}"{{ $tag }}"{{end}}
  ]
}