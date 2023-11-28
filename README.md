# Terraform Provider - Redpanda Brokers

_This provider is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding).

This repository contains a [Terraform](https://www.terraform.io) provider for Redpanda, that enables Terraform to connect to a Redpanda cluster and determine the hostnames of the brokers in use. In provides:

- A single data source (`internal/provider/`),
- Miscellaneous meta files.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

Make sure the provider is listed in required_providers:

```yaml
terraform {
  ...
  required_providers {
    pmwrp = {
      source = "registry.terraform.io/hashicorp/pmwrp"
    }
  ...
  }
}
```

Then, configure the provider by setting the seed address and credentials:

```yaml
provider "pmwrp" {
  seed = "seed-your.cluster.address.here.fmc.prd.cloud.redpanda.com:9092"
  username = "your-username-here"
  password = "your-password-here"
}
```

## Using the data source

Include a data source within your terraform:

```yaml
data "pmwrp_brokers" "default" {}
```

You can then use the broker host names as needed:

```yaml
locals {
  hosts = [for broker in data.pmwrp_brokers.default.brokers: broker.host]
}

output "hosts" {
  value = local.hosts
}
```
