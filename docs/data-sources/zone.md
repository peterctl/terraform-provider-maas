---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "maas_zone Data Source - terraform-provider-maas"
subcategory: ""
description: |-
  Provides details about an existing MAAS zone.
---

# maas_zone (Data Source)

Provides details about an existing MAAS zone.

## Example Usage

```terraform
resource "maas_zone" "test_zone" {
  description = "A description of the test zone"
  name        = "test-zone"
}

data "maas_zone" "test_zone" {
  name = maas_zone.test_zone.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The zone's name.

### Read-Only

- `description` (String) A brief description of the zone.
- `id` (String) The ID of this resource.
