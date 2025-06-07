---
# langfuse_project_api_key Resource

Creates and manages a project API key in Langfuse. Requires an organization scoped API key for authentication.

## Example

```hcl
resource "langfuse_project_api_key" "key" {
  project_id = "proj_123"
  note       = "ci key"
}
```

## Argument Reference

* `project_id` - (Required) ID of the project.
* `note` - (Optional) Optional note for the API key.

## Attributes Reference

* `id` - ID of the API key.
* `public_key` - Public portion of the API key.
* `secret_key` - Secret portion of the API key. Only returned on creation.
* `display_secret_key` - Truncated secret for display.
* `created_at` - Creation timestamp.
* `expires_at` - Expiration timestamp if configured.
* `last_used_at` - Timestamp of last usage if available.

## Import

Project API keys can be imported using the project ID and API key ID separated by a slash:

```shell
terraform import langfuse_project_api_key.key <project_id>/<api_key_id>
```
