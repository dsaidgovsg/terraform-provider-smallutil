provider "smallutil" {
}

data "smallutil_http_req" "plain_tag" {
  url = "http://localhost:8080/plain"
}

data "smallutil_http_req" "json_tag" {
  url = "http://localhost:8080/json"

  response_content_type     = "application/json"
  response_content_json_key = ".tag"
}

data "smallutil_http_req" "failed_tag" {
  url     = "http://localhost:8080/no-such-endpoint"
  default = "failed-tag"
}

data "smallutil_http_req" "override_tag" {
  url      = "http://localhost:8080/no-such-endpoint"
  override = "override-tag"
}

output "plain_tag" {
  value = "${data.smallutil_http_req.plain_tag.value}"
}

output "plain_tag_source" {
  value = "${data.smallutil_http_req.plain_tag.source}"
}

output "json_tag" {
  value = "${data.smallutil_http_req.json_tag.value}"
}

output "json_tag_source" {
  value = "${data.smallutil_http_req.json_tag.source}"
}

output "failed_tag" {
  value = "${data.smallutil_http_req.failed_tag.value}"
}

output "failed_tag_source" {
  value = "${data.smallutil_http_req.failed_tag.source}"
}

output "override_tag" {
  value = "${data.smallutil_http_req.override_tag.value}"
}

output "override_tag_source" {
  value = "${data.smallutil_http_req.override_tag.source}"
}
