---
# langfuse_prompt Resource

Creates and manages a prompt within a Langfuse project.

## Example

```hcl
resource "langfuse_prompt" "example" {
  project_id = "proj_123"
  content    = "Hello, world!"
}
```

## Argument Reference

* `project_id` - (Required) ID of the project that owns the prompt.
* `content` - (Required) Prompt content.

## Attributes Reference

* `id` - Prompt ID.

## Import

Prompts can be imported using `project_id` and `prompt_id` separated by a slash:

```shell
terraform import langfuse_prompt.example <project_id>/<prompt_id>
```
