---
# langfuse_prompt Data Source

Fetches a prompt by name from Langfuse.

## Example

```hcl
data "langfuse_prompt" "example" {
  name = "my_prompt"
}
```

Optional arguments `version` or `label` can be provided to select a specific version.

## Argument Reference

* `name` - (Required) Name of the prompt.
* `version` - (Optional) Version of the prompt to retrieve.
* `label` - (Optional) Deployment label of the prompt.

## Attributes Reference

* `type` - Prompt type (`chat` or `text`).
* `prompt` - JSON encoded prompt content.
* `config` - Configuration map of the prompt.
* `labels` - Deployment labels for the version.
* `tags` - Tags associated with the prompt.
* `commit_message` - Commit message for the version.
* `version_out` - Resolved version number.
