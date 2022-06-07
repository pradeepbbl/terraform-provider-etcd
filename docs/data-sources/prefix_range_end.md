---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "etcd_prefix_range_end Data Source - terraform-provider-etcd"
subcategory: ""
description: |-
  Helper to retrieve a range end that, combined with the key argument, constitutes a prefix of key.
---

# etcd_prefix_range_end (Data Source)

Helper to retrieve a range end that, combined with the key argument, constitutes a prefix of key.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **key** (String) Key to get a prefix of.

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **range_end** (String) Computed range end that, combined with the key, constitutes a prefix of the key.

