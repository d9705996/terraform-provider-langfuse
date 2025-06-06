---
# Langfuse Provider

The Langfuse provider allows Terraform to manage Langfuse projects using the public API.

## Configuration

```hcl
provider "langfuse" {
  host     = "https://cloud.langfuse.com" # optional
  username = "<public key>"
  password = "<secret key>"
}
```

The `username` and `password` can alternatively be provided via the `LANGFUSE_USERNAME` and `LANGFUSE_PASSWORD` environment variables.

## Resources

* [langfuse_project](resources/project.md)
* [langfuse_project_api_key](resources/project_api_key.md)

## Data Sources

* [langfuse_project](data-sources/project.md)
* [langfuse_project_api_keys](data-sources/project_api_keys.md)

