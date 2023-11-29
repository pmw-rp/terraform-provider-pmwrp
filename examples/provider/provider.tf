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
