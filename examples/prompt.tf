provider "langfuse" {
  username = "${var.username}"
  password = "${var.password}"
}

resource "langfuse_prompt" "example" {
  project_id = "proj_123"
  content    = "Hello, world!"
}
