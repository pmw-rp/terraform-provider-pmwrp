---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pmwrp Provider"
subcategory: ""
description: |-
  
---

# pmwrp Provider



## Example Usage

```terraform
terraform {
  required_providers {
    pmwrp = {
      source  = "registry.terraform.io/hashicorp/pmwrp"
      version = "1.0"
    }
  }
}

provider "pmwrp" {
  seed     = "seed-for-your.cluster.address.here:9092"
  username = "your-username-here"
  password = "your-password-here"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `password` (String) password
- `seed` (String) seed broker address
- `username` (String) username
