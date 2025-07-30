---
page_title: "ipnetwork Provider"
description: |-
  IP Network functions
---

# ipnetwork Provider

Provides functions to manipulate IP Network data (address, cidr, ...)

## Example Usage

```terraform
terraform {
  required_providers {
    ipnetwork = {
      source = "jeremmfr/ipnetwork"
    }
  }
  required_version = ">= 1.8.0"
}
```
