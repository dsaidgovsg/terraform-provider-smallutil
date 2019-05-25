# terraform-provider-smallutil

Small utility Terraform provider for command-like data source conveniences.

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
