# terraform-provider-smallutil

[![Build Status](https://travis-ci.org/guangie88/terraform-provider-smallutil.svg?branch=master)](https://travis-ci.org/guangie88/terraform-provider-smallutil)

Small utility Terraform provider for command-like data source conveniences.

Currently the following data sources are supported:

- `smallutil_http_req`
  - Performs a HTTP request (think `curl`) to obtain a value. Currently only
    `text/plain` and `application/json` are supported, and only string value
    from one level of object nesting is supported. This data source allows for
    both `override` value (i.e. if specified, it will just use this value) and
    also `default` value (i.e. if specified, if the HTTP request fails, it will
    fall back to this value).

## How to install

For Linux and Mac OSX users, simply run:

```bash
curl -sL https://raw.githubusercontent.com/guangie88/terraform-provider-smallutil/master/install_from_release.sh | bash
```

This will install the latest version of the plugin and place it into your
Terraform plugin directory (creates the directory if it doesn't exist).

Alternatively, you can visit the
[Releases](https://github.com/guangie88/terraform-provider-smallutil/releases)
page and download the zip file corresponding to your OS and architecture.

Unzip the plugin binary into `~/.terraform.d/plugins`, or
`%APPDATA%\terraform.d\plugins` for Windows. You may check the official
[guide](https://www.terraform.io/docs/plugins/basics.html#installing-plugins)

If you prefer to build the plugin from scratch, follow the step
[here](#how-to-build).

## Terraform `.tf` Example

```tf
provider "smallutil" {}

data "smallutil_http_req" "plain_tag" {
  url = "http://localhost:8080/plain"
}

data "smallutil_http_req" "json_tag" {
  url                       = "http://localhost:8080/json"
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
```

## How to build

You will need Go compiler >= 1.11 (because this use Go modules).

Simply run at the repo root directory:

```bash
go build
```

This will generate `terraform-provider-smallutil` binary, which is the compiled
Terraform provider module.

## Testing

You will need [Terraform](https://www.terraform.io/) CLI.

After building, run the following to move the Terraform provider module to the
test directory:

```bash
mv terraform-provider-smallutil tf_test/
```

After that, you will need to terminals, one to run the test server, and one
to run `terraform` commands.

### Test Server

Navigate to `tf_test/server`. Run the following:

```bash
go build
./server
```

### Test Terraform

Navigate to `tf_test/`. Run the following:

```bash
terraform init
terraform apply
```

You should see that the `terraform apply` command run successfully and shows
a list of Terraform outputs.
