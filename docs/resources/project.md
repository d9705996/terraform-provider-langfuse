---
# langfuse_project Resource

Creates and manages a Langfuse project.

## Example

```hcl
resource "langfuse_project" "example" {
  name      = "demo"
  retention = 7
  metadata = {
    team = "ai"
  }
}
```

## Argument Reference

* `name` - (Required) Name of the project.
* `metadata` - (Optional) Map of metadata values.
* `retention` - (Required) Number of days to retain data. Must be 0 or at least 3.

## Attributes Reference

* `id` - Project ID.
