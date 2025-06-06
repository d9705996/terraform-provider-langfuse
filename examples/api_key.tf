provider "langfuse" {
  username = "${var.username}"
  password = "${var.password}"
}

resource "langfuse_project_api_key" "example" {
  project_id = "proj_123"
  note       = "automation"
}
