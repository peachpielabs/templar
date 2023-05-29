resource "dnsimple_zone_record" "dns_record_{{.subdomain_name | lower}}" {
  zone_name = "${var.dnsimple_domain}"
  name   = "{{.subdomain_name | lower}}"
  value  = "{{.record_value}}"
  type   = "{{.record_type}}"
  ttl    = {{.ttl}}
  url    = "{{.url}}"
}
