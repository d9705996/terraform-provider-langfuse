---
# langfuse_project Data Source

Retrieves information about a Langfuse project by ID.

## Example

```hcl
data "langfuse_project" "example" {
  id = "proj_123"
}
```

## Argument Reference

* `id` - (Required) Project ID.

## Attributes Reference

* `name` - Name of the project.
* `metadata` - Metadata map.
* `retention` - Retention days.
