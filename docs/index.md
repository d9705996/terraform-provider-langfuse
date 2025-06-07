---
# Langfuse Provider

The Langfuse provider allows Terraform to manage Langfuse projects using the public API.

**Note:** This provider is a work in progress. Release artifacts may not yet be available on the GitHub releases page.

## Configuration

```hcl
provider "langfuse" {
  host     = "https://cloud.langfuse.com" # optional
  username = "<public key>"
  password = "<secret key>"
}
```

The `host` can also be set via the `LANGFUSE_HOST` environment variable. The `username` and `password` can alternatively be provided via the `LANGFUSE_USERNAME` and `LANGFUSE_PASSWORD` environment variables.

## Resources

* [langfuse_project](resources/project.md)
* [langfuse_project_api_key](resources/project_api_key.md)
* [langfuse_prompt](resources/prompt.md)

## Data Sources

* [langfuse_project](data-sources/project.md)
* [langfuse_project_api_keys](data-sources/project_api_keys.md)
* [langfuse_prompt](data-sources/prompt.md)

