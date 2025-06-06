provider "langfuse" {
  username = "${var.username}"
  password = "${var.password}"
}

resource "langfuse_project" "example" {
  name      = "example"
  retention = 7
  metadata = {
    env = "test"
  }
}
