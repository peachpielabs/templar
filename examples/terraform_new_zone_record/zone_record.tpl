resource "dnsimple_zone_record" "dns_record_{{.subdomain_name}}" {
  zone_name = "${var.dnsimple_domain}"
  name   = "{{.subdomain_name}}"
  value  = "{{.record_value}}"
  type   = "{{.record_type}}"
  ttl    = {{.ttl}}
  url    = "{{.url}}"
}
