package langfuse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPromptResource(t *testing.T) {
	var deleted bool
	var updated bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			resp := map[string]interface{}{
				"id":      "1",
				"content": req["content"],
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("encode response: %v", err)
			}
		case http.MethodGet:
			content := "hello"
			if updated {
				content = "bye"
			}
			resp := map[string]interface{}{
				"id":      "1",
				"content": content,
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("encode response: %v", err)
			}
		case http.MethodPut:
			updated = true
			var req map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			resp := map[string]interface{}{
				"id":      "1",
				"content": req["content"],
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("encode response: %v", err)
			}
		case http.MethodDelete:
			deleted = true
			w.WriteHeader(204)
		}
	}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(s *terraform.State) error {
			if !deleted {
				return fmt.Errorf("prompt not deleted")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPrompt(server.URL, "hello"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("langfuse_prompt.test", "content", "hello"),
				),
			},
			{
				Config: testAccPrompt(server.URL, "bye"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("langfuse_prompt.test", "content", "bye"),
				),
			},
		},
	})
}

func testAccPrompt(url, content string) string {
	return fmt.Sprintf(`
provider "langfuse" {
  host = "%s"
}

resource "langfuse_prompt" "test" {
  project_id = "123"
  content    = "%s"
}
`, url, content)
}
