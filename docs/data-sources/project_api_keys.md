---
# langfuse_project_api_keys Data Source

Lists API keys for a given project.

## Example

```hcl
data "langfuse_project_api_keys" "all" {
  project_id = "proj_123"
}
```

## Argument Reference

* `project_id` - (Required) ID of the project.

## Attributes Reference

* `api_keys` - List of API keys. Each item contains:
  * `id` - ID of the key.
  * `created_at` - Creation timestamp.
  * `expires_at` - Expiration timestamp if configured.
  * `last_used_at` - Timestamp when the key was last used.
  * `note` - Optional note.
  * `public_key` - Public key.
  * `display_secret_key` - Truncated secret.
